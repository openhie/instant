module.exports = {
  title: 'Instant OpenHIE',
  tagline: 'Simplifying OpenHIE Setup',
  url: 'https://openhie.github.io',
  baseUrl: '/instant/',
  favicon: 'img/favicon.ico',
  organizationName: 'openhie',
  projectName: 'instant',
  themeConfig: {
    algolia: {
      apiKey: '43dfdd6f76217eafc0e68ada109a0251',
      indexName: 'instant',
      algoliaOptions: {} // Optional, if provided by Algolia
    },
    navbar: {
      title: 'Instant OpenHIE',
      logo: {
        alt: 'Instant OpenHIE Logo',
        src: 'img/IOHIE-icon-medium.png'
      },
      items: [
        { to: 'docs/introduction/vision', label: 'Docs', position: 'left' },
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
      copyright: `<a rel="license" href="http://creativecommons.org/licenses/by/4.0/"><img alt="Creative Commons Licence" style="border-width:0" src="https://i.creativecommons.org/l/by/4.0/88x31.png" /></a><br />This work is licensed under a <a rel="license" href="http://creativecommons.org/licenses/by/4.0/">Creative Commons Attribution 4.0 International License</a>.`
    },
    prism: {
      theme: require('prism-react-renderer/themes/nightOwl')
    }
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          editUrl: 'https://github.com/openhie/instant/tree/master/docs/'
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css')
        }
      }
    ]
  ]
}
