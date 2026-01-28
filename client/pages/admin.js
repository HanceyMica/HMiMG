import React, { useEffect, useState } from 'react';
import MainLayout from '../components/MainLayout';
import { Form, Input, Button, Card, message, Select, Row, Col, Tabs } from 'antd';
import api from '../lib/api';
import { useRouter } from 'next/router';
import Cookies from 'js-cookie';
import { useI18n } from '../lib/i18n';

const Admin = (props) => {
  const router = useRouter();
  const [collections, setCollections] = useState([]);
  const [albums, setAlbums] = useState([]);
  const [settingsForm] = Form.useForm();
  const [profileForm] = Form.useForm();
  const { t } = useI18n();

  useEffect(() => {
     // Check if admin
     const userStr = Cookies.get('user');
     if (!userStr) {
         router.push('/login');
         return;
     }
     const user = JSON.parse(userStr);
     if (user.role !== 'admin') {
         message.error('Access denied');
         router.push('/');
         return;
     }
     
     // Set initial profile values
     profileForm.setFieldsValue({
         username: user.username,
         email: user.email, // This might be missing if not in cookie, but we can't fetch full user details easily without new endpoint or re-login. 
         // For now, let's assume cookie has basic info or user will input new values.
         // Better approach: Fetch user details from /auth/me or similar if exists. 
         // Since we don't have /auth/me, we rely on user inputting new values or just what's in cookie.
         // Actually updateAdmin endpoint updates current user.
     });

     fetchData();
  }, []);

  const fetchData = async () => {
    try {
        const [alb, col, sett] = await Promise.all([
            api.get('/albums'),
            api.get('/collections'),
            api.get('/settings')
        ]);
        setAlbums(alb.data);
        setCollections(col.data);
        settingsForm.setFieldsValue(sett.data);
    } catch (e) {}
  };

  const handleUpdateSettings = async (values) => {
      try {
          await api.put('/settings', values);
          message.success(t('admin.settingsUpdated'));
      } catch (e) {
          message.error(t('admin.settingsFailed'));
      }
  };

  const handleUpdateProfile = async (values) => {
      if (values.password) {
          if (values.password !== values.confirmPassword) {
              message.error(t('user.passwordMismatch'));
              return;
          }
          if (!values.oldPassword) {
               // Should be handled by backend, but frontend check is nice
          }
      }

      try {
          const res = await api.put('/admin/update', values);
          
          if (res.data.passwordChanged) {
              message.success(t('user.changePasswordSuccess'));
              // Logout
              Cookies.remove('token');
              Cookies.remove('user');
              setTimeout(() => {
                  router.push('/login');
              }, 1500);
              return;
          }

          message.success(t('user.profileUpdated'));
          // Update cookie if username changed (optional, but good practice)
          const userStr = Cookies.get('user');
          if (userStr) {
              const user = JSON.parse(userStr);
              if (values.username) user.username = values.username;
              Cookies.set('user', JSON.stringify(user));
          }
      } catch (e) {
          message.error(e.response?.data?.error || t('admin.profileFailed'));
      }
  };

  const handleCreateAlbum = async (values) => {
    try {
      await api.post('/albums', values);
      message.success(t('admin.albumCreated'));
      fetchData();
    } catch (err) {
      message.error(t('admin.failedCreateAlbum'));
    }
  };

  const handleCreateCollection = async (values) => {
    try {
      await api.post('/collections', values);
      message.success(t('admin.collectionCreated'));
      fetchData();
    } catch (err) {
      message.error(t('admin.failedCreateCollection'));
    }
  };

  const handleAddToCollection = async (values) => {
      try {
          await api.post('/collections/add', values);
          message.success(t('admin.addedSuccess'));
      } catch (err) {
          message.error(err.response?.data?.message || t('admin.failedAdd'));
      }
  };

  const DashboardContent = () => (
      <>
        <Row gutter={16} style={{ marginBottom: 20 }}>
            <Col span={24}>
                <Card title={t('admin.systemSettings')}>
                    <Form form={settingsForm} onFinish={handleUpdateSettings} layout="vertical">
                        <Row gutter={16}>
                            <Col xs={24} md={8}>
                                <Form.Item name="website_title" label={t('admin.websiteTitle')}>
                                    <Input placeholder={t('admin.websiteTitlePlaceholder')} />
                                </Form.Item>
                            </Col>
                            <Col xs={24} md={8}>
                                <Form.Item name="max_users" label={t('admin.maxUsers')} rules={[{ required: true, message: t('login.required') + t('admin.maxUsers') }]}>
                                    <Input type="number" />
                                </Form.Item>
                            </Col>
                            <Col xs={24} md={8}>
                                <Form.Item name="allow_registration" label={t('admin.allowRegistration')} rules={[{ required: true }]}>
                                    <Select>
                                        <Select.Option value="true">{t('admin.yes')}</Select.Option>
                                        <Select.Option value="false">{t('admin.no')}</Select.Option>
                                    </Select>
                                </Form.Item>
                            </Col>
                        </Row>
                        <Form.Item>
                            <Button type="primary" htmlType="submit">{t('admin.saveSettings')}</Button>
                        </Form.Item>
                    </Form>
                </Card>
            </Col>
        </Row>
        <Row gutter={[16, 16]}>
            <Col xs={24} md={8}>
            <Card title={t('admin.createAlbum')}>
                <Form onFinish={handleCreateAlbum} layout="vertical">
                <Form.Item name="name" rules={[{ required: true, message: t('login.required') + t('admin.name') }]} label={t('admin.name')}>
                    <Input />
                </Form.Item>
                <Form.Item name="description" label={t('admin.description')}>
                    <Input.TextArea />
                </Form.Item>
                <Button type="primary" htmlType="submit">{t('admin.create')}</Button>
                </Form>
            </Card>
            </Col>
            <Col xs={24} md={8}>
            <Card title={t('admin.createCollection')}>
                <Form onFinish={handleCreateCollection} layout="vertical">
                <Form.Item name="name" rules={[{ required: true, message: t('login.required') + t('admin.name') }]} label={t('admin.name')}>
                    <Input />
                </Form.Item>
                <Form.Item name="description" label={t('admin.description')}>
                    <Input.TextArea />
                </Form.Item>
                <Button type="primary" htmlType="submit">{t('admin.create')}</Button>
                </Form>
            </Card>
            </Col>
            <Col xs={24} md={8}>
            <Card title={t('admin.organize')}>
                <Form onFinish={handleAddToCollection} layout="vertical">
                    <Form.Item name="collectionId" label={t('admin.targetCollection')} rules={[{ required: true, message: t('login.required') + t('admin.targetCollection') }]}>
                        <Select>
                            {collections.map(c => <Select.Option key={c.id} value={c.id}>{c.name}</Select.Option>)}
                        </Select>
                    </Form.Item>
                    <Form.Item name="itemType" label={t('admin.itemType')} rules={[{ required: true, message: t('login.required') + t('admin.itemType') }]}>
                        <Select onChange={() => {}}>
                            <Select.Option value="album">{t('admin.album')}</Select.Option>
                            <Select.Option value="collection">{t('admin.collection')}</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item shouldUpdate={(prevValues, currentValues) => prevValues.itemType !== currentValues.itemType}>
                        {({ getFieldValue }) => {
                            const itemType = getFieldValue('itemType');
                            return (
                                <Form.Item 
                                    name="itemName" 
                                    label={t('admin.itemName')} 
                                    rules={[{ required: true, message: t('login.required') + t('admin.itemName') }]}
                                    help={t('admin.itemNameHelp')}
                                >
                                    <Select showSearch optionFilterProp="children">
                                        {itemType === 'album' 
                                            ? albums.map(a => <Select.Option key={a.id} value={a.name}>{a.name}</Select.Option>)
                                            : collections.map(c => <Select.Option key={c.id} value={c.name}>{c.name}</Select.Option>)
                                        }
                                    </Select>
                                </Form.Item>
                            );
                        }}
                    </Form.Item>
                    <Button type="primary" htmlType="submit">{t('admin.add')}</Button>
                </Form>
            </Card>
            </Col>
        </Row>
      </>
  );

  const AccountSettingsContent = () => (
      <Card title={t('admin.updateProfile')} style={{ maxWidth: 600 }}>
          <Form form={profileForm} onFinish={handleUpdateProfile} layout="vertical">
              <Form.Item name="username" label={t('login.username')}>
                  <Input />
              </Form.Item>
              <Form.Item name="email" label={t('login.email')}>
                  <Input type="email" />
              </Form.Item>
              <Form.Item name="phone" label={t('login.phone')}>
                  <Input />
              </Form.Item>
              
              <Form.Item name="oldPassword" label={t('admin.oldPassword')} help={t('admin.passwordHelp')}>
                  <Input.Password />
              </Form.Item>
              
              <Form.Item name="password" label={t('admin.newPassword')}>
                  <Input.Password />
              </Form.Item>
              
              <Form.Item 
                name="confirmPassword" 
                label={t('admin.confirmPassword')}
                dependencies={['password']}
                rules={[
                  ({ getFieldValue }) => ({
                    validator(_, value) {
                      if (!value || getFieldValue('password') === value) {
                        return Promise.resolve();
                      }
                      return Promise.reject(new Error(t('admin.passwordMismatch')));
                    },
                  }),
                ]}
              >
                  <Input.Password />
              </Form.Item>
              
              <Form.Item>
                  <Button type="primary" htmlType="submit">{t('admin.saveSettings')}</Button>
              </Form.Item>
          </Form>
      </Card>
  );

  return (
    <MainLayout {...props}>
      <h1>{t('admin.dashboard')}</h1>
      <Tabs 
        defaultActiveKey="1"
        items={[
            { key: '1', label: t('admin.dashboard'), children: <DashboardContent /> },
            { key: '2', label: t('admin.accountSettings'), children: <AccountSettingsContent /> }
        ]}
      />
    </MainLayout>
  );
};

export default Admin;
