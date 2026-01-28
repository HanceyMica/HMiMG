import React, { useState, useEffect } from 'react';
import { Form, Input, Button, Card, message, theme } from 'antd';
import api from '../lib/api';
import { useRouter } from 'next/router';
import Cookies from 'js-cookie';
import { useI18n } from '../lib/i18n';
import Head from 'next/head';

const LoginContent = ({ isDarkMode }) => {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [isRegister, setIsRegister] = useState(false);
  const [allowRegistration, setAllowRegistration] = useState(false);
  const [websiteTitle, setWebsiteTitle] = useState('HMiMG');
  const { t } = useI18n();
  const { token } = theme.useToken();

  useEffect(() => {
    // Check settings
    api.get('/settings/public')
      .then(res => {
        setAllowRegistration(res.data.allow_registration);
        if (res.data.website_title) {
            setWebsiteTitle(res.data.website_title);
        }
      })
      .catch(err => {
        console.error('Failed to fetch settings', err);
      });
  }, []);

  const onFinish = async (values) => {
    setLoading(true);
    try {
      if (isRegister) {
        await api.post('/register', values);
        message.success(t('login.registerSuccess'));
        setIsRegister(false);
      } else {
        const res = await api.post('/login', values);
        Cookies.set('token', res.data.token);
        Cookies.set('user', JSON.stringify(res.data.user));
        message.success(t('login.success'));
        router.push('/');
      }
    } catch (err) {
      const msg = isRegister ? t('login.registerFailed') : t('login.failed');
      message.error(err.response?.data?.error || msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <Head>
        <title>{isRegister ? t('login.registerBtn') : t('login.title')} - {websiteTitle}</title>
      </Head>
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh', 
        // Background handled by globals.css on body
      }}>
        <Card title={isRegister ? `${t('login.registerBtn')} - ${websiteTitle}` : `${t('login.loginBtn')} - ${websiteTitle}`} style={{ width: 350 }}>
          <Form
            name="auth"
            onFinish={onFinish}
            layout="vertical"
          >
            <Form.Item
              name="username"
              label={t('login.username')}
              rules={[{ required: true, message: t('login.required') + t('login.username') }]}
            >
              <Input placeholder={t('login.username')} />
            </Form.Item>

            <Form.Item
              name="password"
              label={t('login.password')}
              rules={[{ required: true, message: t('login.required') + t('login.password') }]}
            >
              <Input.Password placeholder={t('login.password')} />
            </Form.Item>

            {isRegister && (
              <>
                <Form.Item
                  name="email"
                  label={t('login.email')}
                  rules={[
                    { required: true, message: t('login.required') + t('login.email') },
                    { type: 'email', message: 'Invalid email' }
                  ]}
                >
                  <Input placeholder={t('login.email')} />
                </Form.Item>
                <Form.Item
                  name="phone"
                  label={t('login.phone')}
                  rules={[{ required: true, message: t('login.required') + t('login.phone') }]}
                >
                  <Input placeholder={t('login.phone')} />
                </Form.Item>
              </>
            )}

            <Form.Item>
              <Button type="primary" htmlType="submit" loading={loading} style={{ width: '100%' }}>
                {isRegister ? t('login.registerBtn') : t('login.loginBtn')}
              </Button>
              
              {allowRegistration && (
                <div style={{ marginTop: 10, textAlign: 'center' }}>
                  <Button type="link" onClick={() => setIsRegister(!isRegister)}>
                    {isRegister ? t('login.toLogin') : t('login.toRegister')}
                  </Button>
                </div>
              )}
            </Form.Item>
          </Form>
        </Card>
      </div>
    </>
  );
};

const Login = (props) => {
    return <LoginContent {...props} />
}

export default Login;
