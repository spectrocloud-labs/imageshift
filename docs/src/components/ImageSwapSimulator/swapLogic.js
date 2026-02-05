/**
 * JavaScript port of the Go SwapImage function
 * Based on /api/v1/imageshift_types.go
 */

/**
 * Parse a container image reference into its components
 * @param {string} image - The image reference string
 * @param {string} defaultRegistry - Default registry if none specified
 * @returns {object} Parsed reference object
 */
export function parseImageReference(image, defaultRegistry = 'docker.io') {
  let registry = defaultRegistry;
  let repository = '';
  let tag = 'latest';
  let digest = '';

  // Check for digest
  let imageWithoutDigest = image;
  const digestMatch = image.match(/@(sha256:[a-f0-9]{64})$/);
  if (digestMatch) {
    digest = digestMatch[1];
    imageWithoutDigest = image.replace(/@sha256:[a-f0-9]{64}$/, '');
  }

  // Check for tag
  const tagMatch = imageWithoutDigest.match(/:([^:/]+)$/);
  if (tagMatch) {
    tag = tagMatch[1];
    imageWithoutDigest = imageWithoutDigest.replace(/:[^:/]+$/, '');
  }

  // Parse registry and repository
  const parts = imageWithoutDigest.split('/');

  // Determine if first part is a registry
  if (parts.length === 1) {
    // Just image name, e.g., "nginx"
    repository = `library/${parts[0]}`;
  } else if (parts.length === 2) {
    // Could be user/image or registry/image
    if (parts[0].includes('.') || parts[0].includes(':') || parts[0] === 'localhost') {
      registry = parts[0];
      repository = `library/${parts[1]}`;
    } else {
      repository = `${parts[0]}/${parts[1]}`;
    }
  } else {
    // First part is likely a registry
    if (parts[0].includes('.') || parts[0].includes(':') || parts[0] === 'localhost') {
      registry = parts[0];
      repository = parts.slice(1).join('/');
    } else {
      repository = parts.join('/');
    }
  }

  return {
    registry,
    repository,
    tag,
    digest,
    original: image,
  };
}

/**
 * Format a parsed reference back to a string
 * @param {object} ref - Parsed reference object
 * @returns {string} Formatted image reference
 */
export function formatReference(ref) {
  let result = `${ref.registry}/${ref.repository}`;

  if (ref.tag && ref.tag !== 'latest') {
    result += `:${ref.tag}`;
  } else if (ref.tag === 'latest' && !ref.digest) {
    result += ':latest';
  }

  if (ref.digest) {
    result += `@${ref.digest}`;
  }

  return result;
}

/**
 * Get the full canonical reference string
 * @param {object} ref - Parsed reference object
 * @returns {string} Full canonical reference
 */
export function getCanonicalReference(ref) {
  let result = `${ref.registry}/${ref.repository}`;

  if (ref.digest) {
    // If there's a digest, include both tag and digest if tag exists
    if (ref.tag && ref.tag !== 'latest') {
      result += `:${ref.tag}@${ref.digest}`;
    } else {
      result += `@${ref.digest}`;
    }
  } else {
    result += `:${ref.tag || 'latest'}`;
  }

  return result;
}

/**
 * Swap an image according to the provided mappings
 * @param {string} image - Input image reference
 * @param {object} config - Configuration object with mappings
 * @returns {object} Result with newImage and matchInfo
 */
export function swapImage(image, config) {
  const { defaultRegistry = 'docker.io', mappings = {} } = config;
  const { swap = [], exactSwap = [], regexSwap = [] } = mappings;

  const ref = parseImageReference(image, defaultRegistry);
  let newImage = '';
  let matchType = null;
  let matchRule = null;

  // If registry matches default, construct canonical name
  if (ref.registry === defaultRegistry) {
    newImage = getCanonicalReference(ref);
  }

  // Apply swap rules (registry-level)
  for (const rule of swap) {
    if (rule.registry === ref.registry) {
      const identifier = ref.digest || ref.tag || 'latest';

      if (ref.digest && ref.tag && ref.tag !== 'latest') {
        // Both tag and digest present
        newImage = `${rule.target}/${ref.repository}:${ref.tag}@${ref.digest}`;
      } else if (ref.digest) {
        // Only digest
        newImage = `${rule.target}/${ref.repository}@${ref.digest}`;
      } else {
        // Tag only (or default 'latest')
        newImage = `${rule.target}/${ref.repository}:${identifier}`;
      }

      matchType = 'swap';
      matchRule = rule;
      break;
    }
  }

  // Apply exactSwap rules (can override swap)
  const canonicalRef = getCanonicalReference(ref);
  for (const rule of exactSwap) {
    // Try matching against both the original and canonical form
    if (rule.reference === ref.original || rule.reference === canonicalRef) {
      newImage = rule.target;
      matchType = 'exactSwap';
      matchRule = rule;
      break;
    }
  }

  // Apply regexSwap rules (highest priority)
  for (const rule of regexSwap) {
    try {
      const re = new RegExp(rule.expression);
      const match = canonicalRef.match(re);

      if (match) {
        let result = rule.target;
        const remainingAfterMatch = canonicalRef.replace(match[0], '');

        if (match.length > 1) {
          // Replace capture groups
          for (let m = 1; m < match.length; m++) {
            result = result.replace(`$${m}`, match[m]);
            result = `${result}${remainingAfterMatch}`;
          }
        }

        newImage = result;
        matchType = 'regexSwap';
        matchRule = rule;
        break;
      }
    } catch (e) {
      // Invalid regex, skip this rule
      console.warn(`Invalid regex in regexSwap rule: ${rule.expression}`);
    }
  }

  return {
    original: image,
    parsed: ref,
    newImage: newImage || null,
    matchType,
    matchRule,
  };
}

/**
 * Preset examples for testing
 */
export const presetExamples = [
  {
    name: 'Docker Hub to Internal Registry',
    image: 'nginx:latest',
    config: {
      defaultRegistry: 'docker.io',
      mappings: {
        swap: [
          { registry: 'docker.io', target: 'registry.internal.example.com/dockerhub' },
        ],
      },
    },
  },
  {
    name: 'Multiple Registry Swaps',
    image: 'gcr.io/google-containers/pause:3.2',
    config: {
      defaultRegistry: 'docker.io',
      mappings: {
        swap: [
          { registry: 'docker.io', target: 'internal.example.com/docker' },
          { registry: 'gcr.io', target: 'internal.example.com/gcr' },
          { registry: 'quay.io', target: 'internal.example.com/quay' },
        ],
      },
    },
  },
  {
    name: 'Exact Swap Override',
    image: 'nginx:latest',
    config: {
      defaultRegistry: 'docker.io',
      mappings: {
        swap: [
          { registry: 'docker.io', target: 'internal.example.com/docker' },
        ],
        exactSwap: [
          { reference: 'docker.io/library/nginx:latest', target: 'internal.example.com/approved/nginx:1.25-alpine' },
        ],
      },
    },
  },
  {
    name: 'AWS ECR Regex',
    image: '123456789012.dkr.ecr.us-west-2.amazonaws.com/my-app:v1.0.0',
    config: {
      defaultRegistry: 'docker.io',
      mappings: {
        regexSwap: [
          {
            expression: '^(\\d+)\\.dkr\\.ecr\\.([^.]+)\\.amazonaws\\.com',
            target: 'ecr-mirror.internal.example.com/$1/$2',
          },
        ],
      },
    },
  },
  {
    name: 'Image with Digest',
    image: 'nginx@sha256:abc123def456abc123def456abc123def456abc123def456abc123def456abcd',
    config: {
      defaultRegistry: 'docker.io',
      mappings: {
        swap: [
          { registry: 'docker.io', target: 'internal.example.com/docker' },
        ],
      },
    },
  },
];
