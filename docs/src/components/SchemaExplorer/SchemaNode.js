import React, { useState, useEffect } from 'react';
import styles from './SchemaExplorer.module.css';

function TypeBadge({ type }) {
  const badgeClass = {
    string: styles.typeString,
    object: styles.typeObject,
    array: styles.typeArray,
    boolean: styles.typeBoolean,
    integer: styles.typeInteger,
  }[type] || styles.typeDefault;

  return <span className={`${styles.typeBadge} ${badgeClass}`}>{type}</span>;
}

function RequiredBadge({ required }) {
  return (
    <span className={required ? styles.requiredBadge : styles.optionalBadge}>
      {required ? 'required' : 'optional'}
    </span>
  );
}

export default function SchemaNode({ node, depth = 0, expandAll, defaultExpanded = false }) {
  const [isExpanded, setIsExpanded] = useState(defaultExpanded || depth < 2);
  const hasChildren = node.children && node.children.length > 0;
  const hasItems = node.items && node.items.children && node.items.children.length > 0;
  const isExpandable = hasChildren || hasItems;

  useEffect(() => {
    if (expandAll !== undefined) {
      setIsExpanded(expandAll);
    }
  }, [expandAll]);

  const toggleExpand = () => {
    if (isExpandable) {
      setIsExpanded(!isExpanded);
    }
  };

  const childNodes = hasChildren ? node.children : hasItems ? node.items.children : [];

  return (
    <div className={styles.nodeContainer} style={{ marginLeft: depth > 0 ? '20px' : '0' }}>
      <div className={styles.nodeHeader} onClick={toggleExpand}>
        {isExpandable && (
          <span className={`${styles.expandIcon} ${isExpanded ? styles.expanded : ''}`}>
            {isExpanded ? '▼' : '▶'}
          </span>
        )}
        {!isExpandable && <span className={styles.expandIconPlaceholder} />}

        <span className={styles.nodeName}>{node.name}</span>
        <TypeBadge type={node.type} />
        {node.required !== undefined && <RequiredBadge required={node.required} />}

        {node.default !== undefined && (
          <span className={styles.defaultValue}>
            default: <code>{JSON.stringify(node.default)}</code>
          </span>
        )}
      </div>

      {node.description && (
        <div className={styles.nodeDescription} style={{ marginLeft: isExpandable ? '20px' : '16px' }}>
          {node.description}
        </div>
      )}

      {isExpanded && isExpandable && (
        <div className={styles.childrenContainer}>
          {hasItems && (
            <div className={styles.arrayItemsLabel}>Array items:</div>
          )}
          {childNodes.map((child, index) => (
            <SchemaNode
              key={`${child.name}-${index}`}
              node={child}
              depth={depth + 1}
              expandAll={expandAll}
            />
          ))}
        </div>
      )}
    </div>
  );
}
