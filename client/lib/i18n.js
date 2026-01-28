import React, { createContext, useContext, useState, useEffect } from 'react';
import { en, zh, ja } from '../locales';
import Cookies from 'js-cookie';

const I18nContext = createContext();

export const I18nProvider = ({ children }) => {
  const [locale, setLocale] = useState('zh'); // Default to Chinese as per request implies Chinese user

  useEffect(() => {
    const savedLocale = Cookies.get('locale');
    if (savedLocale) {
      setLocale(savedLocale);
    }
  }, []);

  const changeLocale = (lang) => {
    setLocale(lang);
    Cookies.set('locale', lang);
  };

  const t = (key) => {
    const keys = key.split('.');
    let localeData;
    if (locale === 'en') localeData = en;
    else if (locale === 'ja') localeData = ja;
    else localeData = zh;

    let value = localeData;
    for (const k of keys) {
      value = value?.[k];
    }
    return value || key;
  };

  return (
    <I18nContext.Provider value={{ locale, changeLocale, t }}>
      {children}
    </I18nContext.Provider>
  );
};

export const useI18n = () => useContext(I18nContext);
