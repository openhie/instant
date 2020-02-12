import React from 'react'
import classnames from 'classnames'
import Layout from '@theme/Layout'
import Link from '@docusaurus/Link'
import useDocusaurusContext from '@docusaurus/useDocusaurusContext'
import useBaseUrl from '@docusaurus/useBaseUrl'
import styles from './styles.module.css'

const features = [
  {
    title: <>Introduction</>,
    imageUrl: 'img/undraw_docusaurus_mountain.svg',
    description: (
      <>
        The Instant OpenHIE project aims to reduce the costs and skills required 
        for software developers to deploy an OpenHIE architecture for quicker 
        initial solution testing and as a starting point for faster production 
        implementation and customization.
        <br /><br />
        View the <a href="docs/introduction/overview">Introduction section</a> to learn more.
      </>
    )
  },
  {
    title: <>Concepts</>,
    imageUrl: 'img/undraw_docusaurus_react.svg',
    description: (
      <>
        Instant OpenHIE provides an easy way to setup, explore and develop with
        the OpenHIE Architecture. It allows packages to be added that support
        multiple different use cases and workflows specified by OpenHIE. Each
        package contains scripts to stand up and configure applications that
        support these various workflows.
        <br /><br />
        View the <a href="docs/concepts/overview">Concepts section</a> to learn more.
      </>
    )
  },
  {
    title: <>Packages</>,
    imageUrl: 'img/undraw_docusaurus_tree.svg',
    description: (
      <>
        Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod
        tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim
        veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea
        commodo consequat. Duis aute irure dolor in reprehenderit in voluptate
        velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint
        occaecat cupidatat non proident, sunt in culpa qui officia deserunt
        mollit anim id est laborum
      </>
    )
  }
]

function Feature({ imageUrl, title, description }) {
  const imgUrl = useBaseUrl(imageUrl)
  return (
    <div className={classnames('col col--4', styles.feature)}>
      {imgUrl && (
        <div className="text--center">
          <img className={styles.featureImage} src={imgUrl} alt={title} />
        </div>
      )}
      <h3>{title}</h3>
      <p>{description}</p>
    </div>
  )
}

function Home() {
  const context = useDocusaurusContext()
  const { siteConfig = {} } = context
  return (
    <Layout
      title={`${siteConfig.title}`}
      description="Description will go into a meta tag in <head />"
    >
      <header className={classnames('hero hero--primary', styles.heroBanner)}>
        <div className="container">
          <h1 className="hero__title">{siteConfig.title}</h1>
          <p className="hero__subtitle">{siteConfig.tagline}</p>
          <div className={styles.buttons}>
            <Link
              className={classnames(
                'button button--outline button--secondary button--lg',
                styles.getStarted
              )}
              to={useBaseUrl('docs/introduction/overview')}
            >
              Get Started
            </Link>
          </div>
        </div>
      </header>
      <main>
        {features && features.length && (
          <section className={styles.features}>
            <div className="container">
              <div className="row">
                {features.map((props, idx) => (
                  <Feature key={idx} {...props} />
                ))}
              </div>
            </div>
          </section>
        )}
      </main>
    </Layout>
  )
}

export default Home
