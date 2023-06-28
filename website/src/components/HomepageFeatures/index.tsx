import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';
import Link from "@docusaurus/Link";

type FeatureItem = {
  title: string;
  Svg: React.ComponentType<React.ComponentProps<'svg'>>;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Easy to Use',
    Svg: require('@site/static/img/features/hourglass.svg').default,
    description: (
      <>
        Auditum is a simple pluggable solution for audit trail that integrates
        well with any application with as little operational overhead as
        possible.
      </>
    ),
  },
  {
    title: 'Focus on Your Data',
    Svg: require('@site/static/img/features/data.svg').default,
    description: (
      <>
        Submit audit records and query them with a simple yet powerful API.
        Auditum handles all the routine.
      </>
    )
  },
  {
    title: 'Developer-friendly API',
    Svg: require('@site/static/img/features/api.svg').default,
    description: (
      <>
        Auditum API provides well-documented Protobuf contracts and
        supports HTTP and gRPC protocols. OpenAPI (Swagger) specification
        is also available.
      </>
    )
  },
  {
    title: 'Cloud Native App',
    Svg: require('@site/static/img/features/cloud.svg').default,
    description: (
      <>
        Auditum is a <Link to="https://12factor.net/">12-factor application</Link>,
        it can be run in a container and easily deployed to a Kubernetes cluster.
      </>
    )
  },
  {
    title: 'Observability out of the box',
    Svg: require('@site/static/img/features/telescope.svg').default,
    description: (
      <>
        Auditum has built-in support for structured logging, Prometheus metrics
        and OpenTelemetry tracing.
      </>
    )
  },
  {
    title: 'Open Source',
    Svg: require('@site/static/img/features/opensource.svg').default,
    description: (
      <>
        Powered by Community, Auditum is open source and free. We are open to
        any contributions.
      </>
    )
  }
];

function Feature({title, Svg, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div>
      <div className="text--center padding-horiz--md">
        <h3 className={styles.fontQuantico}>{title}</h3>
        <p className="text--left">{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
