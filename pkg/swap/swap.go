package swap

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	imageshiftv1 "github.com/spectrocloud-labs/imageshift/api/v1"
)

func SwapImage(config imageshiftv1.Imageshift, image string) string {
	ref, _ := name.ParseReference(image, name.WithDefaultRegistry(config.Spec.Default))

	registry := ref.Context().RegistryStr()

	// if registry == default registry return image
	var newImage string

	if registry == config.Spec.Default {
		newImage = ref.Name()
	}

	for _, swap := range config.Spec.Mappings.Swap {
		if swap.Registry == registry {
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
		re := regexp.MustCompile(regexSwap.Expression)

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
