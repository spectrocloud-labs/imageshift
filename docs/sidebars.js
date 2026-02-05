/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  docs: [
    'intro',
    {
      type: 'category',
      label: 'Installation',
      items: [
        'installation/prerequisites',
        'installation/quick-start',
        'installation/helm',
        'installation/kustomize',
      ],
    },
    {
      type: 'category',
      label: 'Reference',
      items: [
        'reference/crd',
        'reference/interactive',
        'reference/configuration',
      ],
    },
  ],
};

export default sidebars;
