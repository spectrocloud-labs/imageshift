import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import CodeBlock from '@theme/CodeBlock';

import styles from './index.module.css';

function HeroSection() {
  const { siteConfig } = useDocusaurusContext();
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <h1 className="hero__title">{siteConfig.title}</h1>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.buttons}>
          <Link
            className="button button--secondary button--lg"
            to="/docs/intro">
            Get Started
          </Link>
          <Link
            className="button button--secondary button--lg"
            href="https://github.com/wcrum/imageshift"
            style={{ marginLeft: '1rem' }}>
            GitHub
          </Link>
        </div>
      </div>
    </header>
  );
}

const features = [
  {
    title: 'Zero Manifest Changes',
    description:
      'Transparently redirect container images without modifying your existing Kubernetes manifests. Just label your namespace and ImageShift handles the rest.',
  },
  {
    title: 'Air-Gapped Ready',
    description:
      'Perfect for disconnected environments. Mirror images to your internal registry and let ImageShift automatically redirect all pulls.',
  },
  {
    title: 'Flexible Mapping',
    description:
      'Choose from registry-level swaps, exact image matches, or powerful regex patterns to handle any image redirection scenario.',
  },
  {
    title: 'Namespace Scoped',
    description:
      'Control which namespaces are affected with simple labels. Enable or disable image swapping on a per-namespace basis.',
  },
];

function FeatureSection() {
  return (
    <section className={styles.features}>
      <div className="container">
        <h2 className={styles.sectionTitle}>Why ImageShift?</h2>
        <div className="row">
          {features.map((feature, idx) => (
            <div key={idx} className={clsx('col col--6', styles.feature)}>
              <div className={styles.featureCard}>
                <h3>{feature.title}</h3>
                <p>{feature.description}</p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

const useCases = [
  {
    title: 'Air-Gapped Environments',
    description: 'Redirect all image pulls to your internal registry in disconnected networks.',
  },
  {
    title: 'Multi-Region Deployments',
    description: 'Use region-specific registries to reduce latency and egress costs.',
  },
  {
    title: 'Image Mirroring',
    description: 'Implement transparent caching and mirroring strategies.',
  },
  {
    title: 'Development & Testing',
    description: 'Test with local registries without changing application manifests.',
  },
];

function UseCasesSection() {
  return (
    <section className={styles.useCases}>
      <div className="container">
        <h2 className={styles.sectionTitle}>Use Cases</h2>
        <div className="row">
          {useCases.map((useCase, idx) => (
            <div key={idx} className={clsx('col col--3', styles.useCase)}>
              <h3>{useCase.title}</h3>
              <p>{useCase.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

const quickStartYaml = `# 1. Label your namespace
kubectl label namespace default imageshift.dev=enabled

# 2. Create an Imageshift resource
cat <<EOF | kubectl apply -f -
apiVersion: imageshift.dev/v1
kind: Imageshift
metadata:
  name: imageshift
spec:
  mappings:
    swap:
      - registry: docker.io
        target: registry.internal.example.com/dockerhub
      - registry: ghcr.io
        target: registry.internal.example.com/ghcr
EOF

# 3. Deploy your workloads - images are automatically redirected!
kubectl run nginx --image=nginx:latest
# Image becomes: registry.internal.example.com/dockerhub/library/nginx:latest`;

function QuickStartSection() {
  return (
    <section className={styles.quickStart}>
      <div className="container">
        <h2 className={styles.sectionTitle}>Quick Start</h2>
        <p className={styles.sectionDescription}>
          Get up and running in minutes with just a few commands.
        </p>
        <CodeBlock language="bash">{quickStartYaml}</CodeBlock>
        <div className={styles.ctaContainer}>
          <Link
            className="button button--primary button--lg"
            to="/docs/installation/quick-start">
            Full Installation Guide
          </Link>
        </div>
      </div>
    </section>
  );
}

export default function Home() {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout
      title={`${siteConfig.title} - Kubernetes Image Redirection`}
      description="Transparent Container Image Redirection for Kubernetes. Automatically swap container images based on configurable mapping rules.">
      <HeroSection />
      <main>
        <FeatureSection />
        <UseCasesSection />
        <QuickStartSection />
      </main>
    </Layout>
  );
}
