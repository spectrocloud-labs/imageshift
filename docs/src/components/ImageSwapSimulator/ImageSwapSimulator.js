import React, { useState, useEffect } from 'react';
import MappingEditor from './MappingEditor';
import SwapResult from './SwapResult';
import { swapImage, presetExamples } from './swapLogic';
import styles from './ImageSwapSimulator.module.css';

const defaultConfig = {
  defaultRegistry: 'docker.io',
  mappings: {
    swap: [{ registry: 'docker.io', target: 'internal.example.com' }],
    exactSwap: [],
    regexSwap: [],
  },
};

function generateYaml(config) {
  const { defaultRegistry, mappings } = config;
  const { swap = [], exactSwap = [], regexSwap = [] } = mappings;

  let yaml = `apiVersion: imageshift.dev/v1
kind: Imageshift
metadata:
  name: imageshift
spec:
  default: ${defaultRegistry}
  namespaceSelector: imageshift.dev`;

  const hasAnyMappings = swap.length > 0 || exactSwap.length > 0 || regexSwap.length > 0;

  if (hasAnyMappings) {
    yaml += `\n  mappings:`;

    if (swap.length > 0) {
      yaml += `\n    swap:`;
      for (const rule of swap) {
        if (rule.registry && rule.target) {
          yaml += `\n      - registry: ${rule.registry}`;
          yaml += `\n        target: ${rule.target}`;
        }
      }
    }

    if (exactSwap.length > 0) {
      yaml += `\n    exactSwap:`;
      for (const rule of exactSwap) {
        if (rule.reference && rule.target) {
          yaml += `\n      - reference: ${rule.reference}`;
          yaml += `\n        target: ${rule.target}`;
        }
      }
    }

    if (regexSwap.length > 0) {
      yaml += `\n    regexSwap:`;
      for (const rule of regexSwap) {
        if (rule.expression && rule.target) {
          yaml += `\n      - expression: "${rule.expression.replace(/\\/g, '\\\\')}"`;
          yaml += `\n        target: "${rule.target}"`;
        }
      }
    }
  }

  return yaml;
}

export default function ImageSwapSimulator() {
  const [image, setImage] = useState('nginx:latest');
  const [config, setConfig] = useState(defaultConfig);
  const [result, setResult] = useState(null);
  const [copied, setCopied] = useState(false);
  const [showYaml, setShowYaml] = useState(false);

  useEffect(() => {
    if (image.trim()) {
      const swapResult = swapImage(image, config);
      setResult(swapResult);
    } else {
      setResult(null);
    }
  }, [image, config]);

  const handleMappingsChange = (mappings) => {
    setConfig({ ...config, mappings });
  };

  const handlePresetChange = (e) => {
    const presetIndex = e.target.value;
    if (presetIndex === '') return;

    const preset = presetExamples[parseInt(presetIndex, 10)];
    setImage(preset.image);
    setConfig(preset.config);
  };

  const handleReset = () => {
    setImage('nginx:latest');
    setConfig(defaultConfig);
  };

  const yaml = generateYaml(config);

  const handleCopyYaml = async () => {
    try {
      await navigator.clipboard.writeText(yaml);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.inputSection}>
        <div className={styles.inputRow}>
          <div className={styles.inputGroup}>
            <label className={styles.label}>Source Image</label>
            <input
              type="text"
              value={image}
              onChange={(e) => setImage(e.target.value)}
              placeholder="Enter image reference (e.g., nginx:latest)"
              className={styles.imageInput}
            />
          </div>

          <div className={styles.inputGroup}>
            <label className={styles.label}>Default Registry</label>
            <input
              type="text"
              value={config.defaultRegistry}
              onChange={(e) => setConfig({ ...config, defaultRegistry: e.target.value })}
              placeholder="Default registry"
              className={styles.registryInput}
            />
          </div>
        </div>

        <div className={styles.presetRow}>
          <select onChange={handlePresetChange} className={styles.presetSelect} defaultValue="">
            <option value="">Load preset example...</option>
            {presetExamples.map((preset, index) => (
              <option key={index} value={index}>
                {preset.name}
              </option>
            ))}
          </select>
          <button onClick={handleReset} className={styles.resetButton}>
            Reset
          </button>
        </div>
      </div>

      <div className={styles.editorSection}>
        <h4 className={styles.sectionTitle}>Mapping Rules</h4>
        <MappingEditor
          mappings={config.mappings}
          onChange={handleMappingsChange}
        />
      </div>

      <SwapResult result={result} />

      <div className={styles.yamlSection}>
        <div className={styles.yamlHeader}>
          <button
            className={styles.yamlToggle}
            onClick={() => setShowYaml(!showYaml)}
          >
            {showYaml ? '▼' : '▶'} Generated YAML
          </button>
          {showYaml && (
            <button className={styles.copyButton} onClick={handleCopyYaml}>
              {copied ? 'Copied!' : 'Copy'}
            </button>
          )}
        </div>
        {showYaml && (
          <pre className={styles.yamlCode}><code>{yaml}</code></pre>
        )}
      </div>
    </div>
  );
}
