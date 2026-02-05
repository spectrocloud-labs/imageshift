import React from 'react';
import styles from './ImageSwapSimulator.module.css';

export default function SwapResult({ result }) {
  if (!result) {
    return (
      <div className={styles.resultContainer}>
        <div className={styles.emptyResult}>
          Enter an image and configure mappings to see the result
        </div>
      </div>
    );
  }

  const { original, parsed, newImage, matchType, matchRule } = result;

  const matchTypeLabels = {
    swap: 'Registry Swap',
    exactSwap: 'Exact Swap',
    regexSwap: 'Regex Swap',
  };

  const matchTypeColors = {
    swap: styles.matchSwap,
    exactSwap: styles.matchExact,
    regexSwap: styles.matchRegex,
  };

  return (
    <div className={styles.resultContainer}>
      <div className={styles.resultHeader}>
        <h4>Result</h4>
        {matchType && (
          <span className={`${styles.matchBadge} ${matchTypeColors[matchType]}`}>
            {matchTypeLabels[matchType]}
          </span>
        )}
      </div>

      <div className={styles.resultGrid}>
        <div className={styles.resultRow}>
          <span className={styles.resultLabel}>Input:</span>
          <code className={styles.resultValue}>{original}</code>
        </div>

        <div className={styles.resultRow}>
          <span className={styles.resultLabel}>Parsed:</span>
          <div className={styles.parsedDetails}>
            <span><strong>Registry:</strong> {parsed.registry}</span>
            <span><strong>Repository:</strong> {parsed.repository}</span>
            <span><strong>Tag:</strong> {parsed.tag}</span>
            {parsed.digest && <span><strong>Digest:</strong> {parsed.digest}</span>}
          </div>
        </div>

        <div className={`${styles.resultRow} ${styles.outputRow}`}>
          <span className={styles.resultLabel}>Output:</span>
          {newImage ? (
            <code className={styles.resultValueOutput}>{newImage}</code>
          ) : (
            <span className={styles.noMatch}>No matching rule - image unchanged</span>
          )}
        </div>

        {matchRule && (
          <div className={styles.resultRow}>
            <span className={styles.resultLabel}>Matched Rule:</span>
            <code className={styles.matchedRule}>
              {matchType === 'swap' && `${matchRule.registry} → ${matchRule.target}`}
              {matchType === 'exactSwap' && `${matchRule.reference} → ${matchRule.target}`}
              {matchType === 'regexSwap' && `/${matchRule.expression}/ → ${matchRule.target}`}
            </code>
          </div>
        )}
      </div>
    </div>
  );
}
