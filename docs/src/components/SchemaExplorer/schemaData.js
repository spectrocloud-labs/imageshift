/**
 * CRD Schema Definition
 * Based on /api/v1/imageshift_types.go
 */

export const schemaData = {
  name: 'Imageshift',
  type: 'object',
  description: 'Imageshift is the Schema for the imageshifts API',
  children: [
    {
      name: 'apiVersion',
      type: 'string',
      required: true,
      description: 'API version of the resource',
      default: 'imageshift.dev/v1',
    },
    {
      name: 'kind',
      type: 'string',
      required: true,
      description: 'Kind of the resource',
      default: 'Imageshift',
    },
    {
      name: 'metadata',
      type: 'object',
      required: true,
      description: 'Standard Kubernetes metadata',
      children: [
        {
          name: 'name',
          type: 'string',
          required: true,
          description: 'Name of the Imageshift resource',
        },
      ],
    },
    {
      name: 'spec',
      type: 'object',
      required: false,
      description: 'ImageshiftSpec defines the desired state of Imageshift',
      children: [
        {
          name: 'default',
          type: 'string',
          required: false,
          default: 'docker.io',
          description: 'The default registry for images without an explicit registry',
        },
        {
          name: 'namespaceSelector',
          type: 'string',
          required: true,
          default: 'imageshift.dev',
          description: 'The label key used to identify namespaces for image swapping',
        },
        {
          name: 'mappings',
          type: 'object',
          required: false,
          description: 'Container for all mapping rules',
          children: [
            {
              name: 'swap',
              type: 'array',
              required: false,
              description: 'Registry-level swaps - redirect all images from one registry to another',
              items: {
                type: 'object',
                children: [
                  {
                    name: 'registry',
                    type: 'string',
                    required: true,
                    description: 'Source registry to match (e.g., docker.io, gcr.io)',
                  },
                  {
                    name: 'target',
                    type: 'string',
                    required: true,
                    description: 'Target registry to redirect to',
                  },
                ],
              },
            },
            {
              name: 'exactSwap',
              type: 'array',
              required: false,
              description: 'Exact image swaps - replace specific image references',
              items: {
                type: 'object',
                children: [
                  {
                    name: 'reference',
                    type: 'string',
                    required: true,
                    description: 'Exact image reference to match',
                  },
                  {
                    name: 'target',
                    type: 'string',
                    required: true,
                    description: 'Replacement image reference',
                  },
                ],
              },
            },
            {
              name: 'regexSwap',
              type: 'array',
              required: false,
              description: 'Regex swaps - use regular expressions for complex matching',
              items: {
                type: 'object',
                children: [
                  {
                    name: 'expression',
                    type: 'string',
                    required: true,
                    description: 'Regular expression pattern to match',
                  },
                  {
                    name: 'target',
                    type: 'string',
                    required: true,
                    description: 'Replacement string (supports capture groups like $1, $2)',
                  },
                ],
              },
            },
          ],
        },
      ],
    },
    {
      name: 'status',
      type: 'object',
      required: false,
      description: 'ImageshiftStatus defines the observed state of Imageshift (currently unused)',
      children: [],
    },
  ],
};
