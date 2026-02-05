import React from 'react';
import styles from './ImageSwapSimulator.module.css';

function SwapRuleEditor({ rules, onChange }) {
  const addRule = () => {
    onChange([...rules, { registry: '', target: '' }]);
  };

  const updateRule = (index, field, value) => {
    const newRules = [...rules];
    newRules[index] = { ...newRules[index], [field]: value };
    onChange(newRules);
  };

  const removeRule = (index) => {
    onChange(rules.filter((_, i) => i !== index));
  };

  return (
    <div className={styles.ruleEditor}>
      {rules.length === 0 && (
        <div className={styles.emptyState}>No swap rules configured</div>
      )}
      {rules.map((rule, index) => (
        <div key={index} className={styles.ruleRow}>
          <input
            type="text"
            placeholder="Source registry (e.g., docker.io)"
            value={rule.registry}
            onChange={(e) => updateRule(index, 'registry', e.target.value)}
            className={styles.ruleInput}
          />
          <span className={styles.arrow}>→</span>
          <input
            type="text"
            placeholder="Target registry"
            value={rule.target}
            onChange={(e) => updateRule(index, 'target', e.target.value)}
            className={styles.ruleInput}
          />
          <button
            onClick={() => removeRule(index)}
            className={styles.removeButton}
            title="Remove rule"
          >
            ×
          </button>
        </div>
      ))}
      <button onClick={addRule} className={styles.addButton}>
        + Add Swap Rule
      </button>
    </div>
  );
}

function ExactSwapRuleEditor({ rules, onChange }) {
  const addRule = () => {
    onChange([...rules, { reference: '', target: '' }]);
  };

  const updateRule = (index, field, value) => {
    const newRules = [...rules];
    newRules[index] = { ...newRules[index], [field]: value };
    onChange(newRules);
  };

  const removeRule = (index) => {
    onChange(rules.filter((_, i) => i !== index));
  };

  return (
    <div className={styles.ruleEditor}>
      {rules.length === 0 && (
        <div className={styles.emptyState}>No exact swap rules configured</div>
      )}
      {rules.map((rule, index) => (
        <div key={index} className={styles.ruleRow}>
          <input
            type="text"
            placeholder="Exact image reference"
            value={rule.reference}
            onChange={(e) => updateRule(index, 'reference', e.target.value)}
            className={styles.ruleInput}
          />
          <span className={styles.arrow}>→</span>
          <input
            type="text"
            placeholder="Target image"
            value={rule.target}
            onChange={(e) => updateRule(index, 'target', e.target.value)}
            className={styles.ruleInput}
          />
          <button
            onClick={() => removeRule(index)}
            className={styles.removeButton}
            title="Remove rule"
          >
            ×
          </button>
        </div>
      ))}
      <button onClick={addRule} className={styles.addButton}>
        + Add Exact Swap Rule
      </button>
    </div>
  );
}

function RegexSwapRuleEditor({ rules, onChange }) {
  const addRule = () => {
    onChange([...rules, { expression: '', target: '' }]);
  };

  const updateRule = (index, field, value) => {
    const newRules = [...rules];
    newRules[index] = { ...newRules[index], [field]: value };
    onChange(newRules);
  };

  const removeRule = (index) => {
    onChange(rules.filter((_, i) => i !== index));
  };

  return (
    <div className={styles.ruleEditor}>
      {rules.length === 0 && (
        <div className={styles.emptyState}>No regex swap rules configured</div>
      )}
      {rules.map((rule, index) => (
        <div key={index} className={styles.ruleRow}>
          <input
            type="text"
            placeholder="Regex expression"
            value={rule.expression}
            onChange={(e) => updateRule(index, 'expression', e.target.value)}
            className={styles.ruleInput}
          />
          <span className={styles.arrow}>→</span>
          <input
            type="text"
            placeholder="Target (use $1, $2 for groups)"
            value={rule.target}
            onChange={(e) => updateRule(index, 'target', e.target.value)}
            className={styles.ruleInput}
          />
          <button
            onClick={() => removeRule(index)}
            className={styles.removeButton}
            title="Remove rule"
          >
            ×
          </button>
        </div>
      ))}
      <button onClick={addRule} className={styles.addButton}>
        + Add Regex Swap Rule
      </button>
    </div>
  );
}

export default function MappingEditor({ mappings, onChange }) {
  const [activeTab, setActiveTab] = React.useState('swap');

  const updateMappings = (type, rules) => {
    onChange({
      ...mappings,
      [type]: rules,
    });
  };

  const tabs = [
    { id: 'swap', label: 'swap', count: mappings.swap?.length || 0 },
    { id: 'exactSwap', label: 'exactSwap', count: mappings.exactSwap?.length || 0 },
    { id: 'regexSwap', label: 'regexSwap', count: mappings.regexSwap?.length || 0 },
  ];

  return (
    <div className={styles.mappingEditor}>
      <div className={styles.tabs}>
        {tabs.map((tab) => (
          <button
            key={tab.id}
            className={`${styles.tab} ${activeTab === tab.id ? styles.activeTab : ''}`}
            onClick={() => setActiveTab(tab.id)}
          >
            {tab.label}
            {tab.count > 0 && <span className={styles.tabBadge}>{tab.count}</span>}
          </button>
        ))}
      </div>
      <div className={styles.tabContent}>
        {activeTab === 'swap' && (
          <SwapRuleEditor
            rules={mappings.swap || []}
            onChange={(rules) => updateMappings('swap', rules)}
          />
        )}
        {activeTab === 'exactSwap' && (
          <ExactSwapRuleEditor
            rules={mappings.exactSwap || []}
            onChange={(rules) => updateMappings('exactSwap', rules)}
          />
        )}
        {activeTab === 'regexSwap' && (
          <RegexSwapRuleEditor
            rules={mappings.regexSwap || []}
            onChange={(rules) => updateMappings('regexSwap', rules)}
          />
        )}
      </div>
    </div>
  );
}
