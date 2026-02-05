import React, { useState, useEffect } from 'react';

export default function VersionBadge() {
  const [version, setVersion] = useState('');

  useEffect(() => {
    fetch('https://api.github.com/repos/wcrum/imageshift/releases/latest')
      .then((res) => res.json())
      .then((data) => {
        if (data.tag_name) {
          setVersion(data.tag_name);
        }
      })
      .catch(() => {
        // Silently fail - badge just won't show
      });
  }, []);

  if (!version) {
    return null;
  }

  return (
    <a
      href="https://github.com/wcrum/imageshift/releases"
      target="_blank"
      rel="noopener noreferrer"
      className="navbar__link version-badge"
    >
      {version}
    </a>
  );
}
