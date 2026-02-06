package swap

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/google/go-containerregistry/pkg/name"
	imageshiftv1 "github.com/spectrocloud-labs/imageshift/api/v1"
)

// regexCache caches compiled regex patterns to avoid recompilation per pod
var (
	regexCache   = make(map[string]*regexp.Regexp)
	regexCacheMu sync.RWMutex
)

// getCompiledRegex returns a cached compiled regex or compiles and caches a new one
func getCompiledRegex(expression string) (*regexp.Regexp, error) {
	regexCacheMu.RLock()
	if re, ok := regexCache[expression]; ok {
		regexCacheMu.RUnlock()
		return re, nil
	}
	regexCacheMu.RUnlock()

	regexCacheMu.Lock()
	defer regexCacheMu.Unlock()

	// Double-check after acquiring write lock
	if re, ok := regexCache[expression]; ok {
		return re, nil
	}

	re, err := regexp.Compile(expression)
	if err != nil {
		return nil, err
	}
	regexCache[expression] = re
	return re, nil
}

// normalizeRegistry returns the canonical registry name using go-containerregistry
func normalizeRegistry(registry string) string {
	// Parse a dummy image with the registry to get the canonical form
	ref, err := name.ParseReference(registry + "/dummy:latest")
	if err != nil {
		return registry
	}
	return ref.Context().RegistryStr()
}

func SwapImage(config imageshiftv1.Imageshift, image string) string {
	ref, _ := name.ParseReference(image, name.WithDefaultRegistry(config.Spec.Default))

	registry := ref.Context().RegistryStr()
	normalizedDefault := normalizeRegistry(config.Spec.Default)

	// if registry == default registry return image
	var newImage string

	if registry == normalizedDefault {
		newImage = ref.Name()
	}

	for _, swap := range config.Spec.Mappings.Swap {
		normalizedSwapRegistry := normalizeRegistry(swap.Registry)
		if normalizedSwapRegistry == registry {
			identifier := ref.Identifier()
			switch len(strings.Split(identifier, ":")) {
			case 1:
				newImage = fmt.Sprintf("%s/%s:%s", swap.Target, ref.Context().RepositoryStr(), ref.Identifier())
			case 2:
				tag := strings.Split(strings.Split(image, "@")[0], ":")[1]
				newImage = fmt.Sprintf("%s/%s:%s@%s", swap.Target, ref.Context().RepositoryStr(), tag, ref.Identifier())
			default:
				newImage = fmt.Sprintf("%s/%s", swap.Target, ref.Context().RepositoryStr())
			}
			break
		}
	}

	for _, exactSwap := range config.Spec.Mappings.ExactSwap {
		if exactSwap.Reference == ref.String() {
			newImage = exactSwap.Target
			break
		}
	}

	for _, regexSwap := range config.Spec.Mappings.RegexSwap {
		re, err := getCompiledRegex(regexSwap.Expression)
		if err != nil {
			// Skip invalid regex patterns
			continue
		}

		match := re.FindStringSubmatch(ref.String())
		if match != nil {
			newRepository := strings.Replace(ref.String(), match[0], "", -1)

			newImage = regexSwap.Target

			if len(match) > 1 {
				for m := 1; m < len(match); m++ {
					newImage = strings.Replace(newImage, "$"+strconv.Itoa(m), match[m], -1)
					newImage = fmt.Sprintf("%s%s", newImage, newRepository)
				}
			}
			break
		}
	}

	return newImage
}
