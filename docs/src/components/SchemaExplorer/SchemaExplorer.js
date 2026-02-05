import React, { useState } from 'react';
import SchemaNode from './SchemaNode';
import { schemaData } from './schemaData';
import styles from './SchemaExplorer.module.css';

export default function SchemaExplorer() {
  const [expandAll, setExpandAll] = useState(undefined);

  const handleExpandAll = () => {
    setExpandAll(true);
    setTimeout(() => setExpandAll(undefined), 100);
  };

  const handleCollapseAll = () => {
    setExpandAll(false);
    setTimeout(() => setExpandAll(undefined), 100);
  };

  return (
    <div className={styles.container}>
      <div className={styles.toolbar}>
        <button className={styles.toolbarButton} onClick={handleExpandAll}>
          Expand All
        </button>
        <button className={styles.toolbarButton} onClick={handleCollapseAll}>
          Collapse All
        </button>
      </div>
      <div className={styles.treeContainer}>
        <SchemaNode node={schemaData} depth={0} expandAll={expandAll} defaultExpanded={true} />
      </div>
    </div>
  );
}
