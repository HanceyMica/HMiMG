import React from 'react';
import { ConfigProvider, theme } from 'antd';
import 'antd/dist/reset.css';
import '../styles/globals.css'; 
import { I18nProvider } from '../lib/i18n';

// We need to create globals.css too for basic body reset if needed, but antd handles most.

const App = ({ Component, pageProps }) => {
  const [isDarkMode, setIsDarkMode] = React.useState(false);

  React.useEffect(() => {
    // Check system preference
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
      setIsDarkMode(true);
    }

    // Listen for changes
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    const handleChange = (e) => setIsDarkMode(e.matches);
    mediaQuery.addEventListener('change', handleChange);

    return () => mediaQuery.removeEventListener('change', handleChange);
  }, []);

  React.useEffect(() => {
    document.body.className = isDarkMode ? 'dark-mode' : 'light-mode';
  }, [isDarkMode]);

  return (
    <I18nProvider>
      <ConfigProvider
        theme={{
          algorithm: isDarkMode ? theme.darkAlgorithm : theme.defaultAlgorithm,
          token: {
            colorBgContainer: isDarkMode ? 'rgba(0, 0, 0, 0.6)' : 'rgba(255, 255, 255, 0.75)',
            colorBgElevated: isDarkMode ? 'rgba(30, 30, 30, 0.8)' : 'rgba(255, 255, 255, 0.9)',
          },
        }}
      >
        <Component {...pageProps} isDarkMode={isDarkMode} setIsDarkMode={setIsDarkMode} />
      </ConfigProvider>
    </I18nProvider>
  );
};

export default App;
