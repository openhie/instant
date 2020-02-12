module.exports = {
  title: 'InstantHIE',
  tagline: 'Simplifying OpenHIE Setup',
  url: 'https://openhie.github.io',
  baseUrl: '/instant/',
  favicon: 'img/favicon.ico',
  organizationName: 'openhie',
  projectName: 'instant',
  themeConfig: {
    algolia: {
      apiKey: '',
      indexName: 'InstantHIE',
      algoliaOptions: {}, // Optional, if provided by Algolia
    },
    navbar: {
      title: 'InstantHIE',
      logo: {
        alt: 'My Site Logo',
        src: 'img/logo.png'
      },
      links: [
        {to: 'docs/introduction/overview', label: 'Docs', position: 'left'},
        {
          href: 'https://github.com/openhie/instant',
          label: 'GitHub',
          position: 'right'
        }
      ]
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Community',
          items: [
            {
              label: 'OpenHIE',
              href: 'https://ohie.org/'
            }
          ]
        }
      ],
      copyright: `Copyright © ${new Date().getFullYear()} OpenHIE`
    },
    prism: {
      theme: require('prism-react-renderer/themes/nightOwl'),
    }
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          editUrl: 'https://github.com/facebook/docusaurus/edit/master/website/'
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css')
        }
      }
    ]
  ]
}
