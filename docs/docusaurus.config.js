module.exports = {
  title: 'InstantHIE',
  tagline: 'Simplifying OpenHIE Setup',
  url: 'https://your-docusaurus-test-site.com',
  baseUrl: '/',
  favicon: 'img/favicon.ico',
  organizationName: 'OpenHIE',
  projectName: 'Instant OpenHIE',
  themeConfig: {
    navbar: {
      title: 'InstantHIE',
      logo: {
        alt: 'My Site Logo',
        src: 'img/logo.svg'
      },
      links: [
        { to: 'docs/doc1', label: 'Docs', position: 'left' },
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
