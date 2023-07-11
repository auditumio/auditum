// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Auditum',
  tagline: 'Audit Log management for any application. Cloud-native, developer-friendly and open-source.',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'https://auditum.io',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'infragmo', // Usually your GitHub org/user name.
  projectName: 'auditum', // Usually your repo name.

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl:
            'https://github.com/infragmo/auditum/tree/main/website/',
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl:
              'https://github.com/infragmo/auditum/tree/main/website/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
        googleAnalytics: {
          trackingID: 'G-X8HNJ4H823',
        }
      }),
    ],
    [
      'redocusaurus',
      {
        debug: Boolean(process.env.DEBUG),
        specs: [
          {
            spec: 'redocusaurus/openapi/v1alpha1/api.yaml',
            route: '/reference/redoc/latest',
          },
          {
            spec: 'redocusaurus/openapi/v1alpha1/api.yaml',
            route: '/reference/redoc/v1alpha1',
          },
        ],
        theme: {
          primaryColor: '#a528d2',
          primaryColorDark: '#c669e8',
        },
        config: 'redocusaurus/redocly.yaml'
      },
    ],
  ],

  // Reference: https://docusaurus.io/docs/api/themes/configuration
  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      image: 'img/social-card.png',
      colorMode: {
        defaultMode: 'light',
        disableSwitch: false,
      },
      announcementBar: {
        id: 'announcementBar-1', // Increment on change
        content:
          '⭐️ Auditum just launched, and we would appreciate if you spread the word! Give it a star on <a target="_blank" rel="noopener noreferrer" href="https://github.com/infragmo/auditum">GitHub</a> and follow us on <a target="_blank" rel="noopener noreferrer" href="https://twitter.com/auditumio">Twitter</a>',
        backgroundColor: '#fafbfc',
        textColor: '#091E42',
        isCloseable: true,
      },
      navbar: {
        title: 'Auditum',
        logo: {
          alt: 'Auditum Logo',
          src: 'img/logo.svg',
        },
        items: [
          {
            type: 'docSidebar',
            sidebarId: 'docSidebar',
            position: 'left',
            label: 'Docs',
          },
          {
            to: '/blog',
            label: 'Blog',
            position: 'left',
          },
          {
            to: '/reference/redoc/latest',
            label: 'API Reference',
            position: 'left',
          },
          {
            href: 'https://github.com/infragmo/auditum',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },
      algolia: {
        appId: '50103BGEZC',
        apiKey: 'c8ed699ceb86243c34bd7d7089bf8c93',
        indexName: 'auditum',
        contextualSearch: true,
        searchParameters: {},
        searchPagePath: 'search',
      },
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Docs',
            items: [
              {
                label: 'Introduction',
                to: '/docs/intro',
              },
              {
                label: 'Installation',
                to: '/docs/getting-started/installation',
              },
              {
                label: 'Configuration',
                to: '/docs/getting-started/configuration',
              },
              {
                label: 'Usage Guide',
                to: '/docs/usage-guide',
              },
            ],
          },
          {
            title: 'About',
            items: [
              {
                label: 'Blog',
                to: '/blog',
              },
              {
                label: 'Credits',
                to: '/credits',
              }
            ]
          },
          {
            title: 'Community',
            items: [
              {
                label: 'GitHub',
                href: 'https://github.com/infragmo/auditum',
              },
              {
                label: 'Twitter',
                href: 'https://twitter.com/auditumio',
              },
              {
                label: 'Stack Overflow',
                href: 'https://stackoverflow.com/questions/tagged/auditum',
              },
            ],
          },
        ],
        copyright: `Auditum 2023-now • Designed and built with 💜 by <a href="https://infragmo.com/?utm_source=auditum-website&utm_medium=footer" style="color: var(--ifm-color-primary-lightest)">Infragmo</a>`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
    }),
};

module.exports = config;
