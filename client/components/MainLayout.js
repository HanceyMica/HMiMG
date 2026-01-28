import React, { useState, useEffect } from 'react';
import { Layout, Menu, Button, Switch, theme, Select, Drawer, Grid } from 'antd';
import { useRouter } from 'next/router';
import Cookies from 'js-cookie';
import Link from 'next/link';
import Head from 'next/head';
import api from '../lib/api';
import { useI18n } from '../lib/i18n';
import { MenuOutlined, SunOutlined, MoonOutlined } from '@ant-design/icons';

const { Header, Content, Footer } = Layout;
const { useBreakpoint } = Grid;

const MainLayout = ({ children, isDarkMode, setIsDarkMode }) => {
  const router = useRouter();
  const { t, locale, changeLocale } = useI18n();
  const screens = useBreakpoint();
  const [drawerVisible, setDrawerVisible] = useState(false);
  const [websiteTitle, setWebsiteTitle] = useState('HMiMG');
  const {
    token: { colorBgContainer },
  } = theme.useToken();

  useEffect(() => {
    // Fetch website title
    api.get('/settings/public')
      .then(res => {
        if (res.data.website_title) {
            setWebsiteTitle(res.data.website_title);
        }
      })
      .catch(e => console.error(e));
  }, []);

  const getSelectedKey = () => {
      const path = router.pathname;
      if (path === '/') return ['home'];
      if (path.startsWith('/library') || path.startsWith('/album') || path.startsWith('/collection') || path.startsWith('/image')) return ['library'];
      if (path.startsWith('/admin')) return ['admin'];
      if (path.startsWith('/about')) return ['about'];
      return ['home'];
  };

  const getPageTitle = () => {
    const path = router.pathname;
    if (path === '/') return t('common.home');
    if (path.startsWith('/library')) return t('common.library');
    if (path.startsWith('/album')) return t('common.albums');
    if (path.startsWith('/collection')) return t('common.collections');
    if (path.startsWith('/image')) return t('image.details');
    if (path.startsWith('/admin')) return t('common.admin');
    if (path.startsWith('/about')) return t('common.about');
    return t('common.home');
  };

  const handleLogout = () => {
    Cookies.remove('token');
    router.push('/login');
  };

  const menuItems = [
    { key: 'home', label: <Link href="/">{t('common.home')}</Link> },
    { key: 'library', label: <Link href="/library">{t('common.library')}</Link> },
    { key: 'admin', label: <Link href="/admin">{t('common.admin')}</Link> },
    { key: 'about', label: <Link href="/about">{t('common.about')}</Link> },
  ];

  const renderRightContent = () => (
      <div style={{ display: 'flex', alignItems: 'center', gap: '15px', flexDirection: screens.md ? 'row' : 'column' }}>
           <Select 
              value={locale} 
              onChange={changeLocale} 
              style={{ width: 100 }}
              options={[
                  { value: 'en', label: 'English' },
                  { value: 'zh', label: '简体中文' },
                  { value: 'ja', label: '日本語' }
              ]}
           />
           <Switch 
              checkedChildren={<MoonOutlined />} 
              unCheckedChildren={<SunOutlined />} 
              checked={isDarkMode}
              onChange={() => setIsDarkMode(!isDarkMode)}
           />
           <Button type="primary" danger onClick={handleLogout}>{t('common.logout')}</Button>
      </div>
  );

  return (
    <>
      <Head>
        <title>{getPageTitle()} - {websiteTitle}</title>
      </Head>
      <Layout className="layout" style={{ minHeight: '100vh' }}>
        <Header style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: screens.md ? '0 50px' : '0 20px' }}>
          <div className="logo" style={{ color: 'white', fontSize: '1.5rem', fontWeight: 'bold', marginRight: '20px' }}>
            {websiteTitle}
          </div>
          
          {screens.md ? (
              <>
                  <Menu
                    theme="dark"
                    mode="horizontal"
                    selectedKeys={getSelectedKey()}
                    items={menuItems}
                    style={{ flex: 1, minWidth: 0 }}
                  />
                  {renderRightContent()}
              </>
          ) : (
              <>
                  <Button type="primary" icon={<MenuOutlined />} onClick={() => setDrawerVisible(true)} />
                  <Drawer
                      title="Menu"
                      placement="right"
                      onClose={() => setDrawerVisible(false)}
                      open={drawerVisible}
                  >
                      <Menu
                          mode="vertical"
                          selectedKeys={getSelectedKey()}
                          items={menuItems}
                          onClick={() => setDrawerVisible(false)}
                          style={{ borderRight: 0, marginBottom: 20 }}
                      />
                      {renderRightContent()}
                  </Drawer>
              </>
          )}
        </Header>
        <Content style={{ padding: screens.md ? '0 50px' : '0 10px', marginTop: '20px' }}>
          <div style={{ background: colorBgContainer, padding: screens.md ? 24 : 12, minHeight: 280, borderRadius: '8px' }}>
            {children}
          </div>
        </Content>
        <Footer style={{ textAlign: 'center' }}>{t('common.copyright')}</Footer>
      </Layout>
    </>
  );
};

export default MainLayout;
