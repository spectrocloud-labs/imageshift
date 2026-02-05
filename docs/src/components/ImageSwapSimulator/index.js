import React from 'react';
import BrowserOnly from '@docusaurus/BrowserOnly';

export default function ImageSwapSimulator() {
  return (
    <BrowserOnly fallback={<div>Loading swap simulator...</div>}>
      {() => {
        const ImageSwapSimulatorComponent = require('./ImageSwapSimulator').default;
        return <ImageSwapSimulatorComponent />;
      }}
    </BrowserOnly>
  );
}
