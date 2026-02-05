import React from 'react';
import BrowserOnly from '@docusaurus/BrowserOnly';

export default function SchemaExplorer() {
  return (
    <BrowserOnly fallback={<div>Loading schema explorer...</div>}>
      {() => {
        const SchemaExplorerComponent = require('./SchemaExplorer').default;
        return <SchemaExplorerComponent />;
      }}
    </BrowserOnly>
  );
}
