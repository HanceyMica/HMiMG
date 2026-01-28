import React, { useEffect, useState } from 'react';
import MainLayout from '../components/MainLayout';
import { Card, Button, Row, Col, Statistic, Upload, message, Modal, Select, Form } from 'antd';
import { PictureOutlined, FolderOutlined, CloudUploadOutlined, InboxOutlined } from '@ant-design/icons';
import { useRouter } from 'next/router';
import Cookies from 'js-cookie';
import api from '../lib/api';
import { useI18n } from '../lib/i18n';
import Link from 'next/link';

const { Dragger } = Upload;

const Home = (props) => {
  const router = useRouter();
  const { t } = useI18n();
  const [albums, setAlbums] = useState([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [fileList, setFileList] = useState([]);
  const [uploading, setUploading] = useState(false);
  const [form] = Form.useForm();

  useEffect(() => {
    const token = Cookies.get('token');
    if (!token) {
      router.push('/login');
      return;
    }
    fetchAlbums();
  }, []);

  const fetchAlbums = async () => {
      try {
          const res = await api.get('/albums');
          setAlbums(res.data);
      } catch (e) {
          console.error(e);
      }
  };

  const handleBeforeUpload = (file) => {
      setFileList(prev => [...prev, file]);
      if (!isModalOpen) {
          setIsModalOpen(true);
      }
      return false; // Prevent auto upload
  };

  const handleUpload = async (values) => {
      if (fileList.length === 0) return;
      
      const formData = new FormData();
      fileList.forEach(file => {
          formData.append('images', file);
      });
      formData.append('albumId', values.albumId);

      setUploading(true);
      try {
          await api.post('/upload', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
          });
          message.success(t('home.uploadSuccess'));
          setIsModalOpen(false);
          setFileList([]);
          form.resetFields();
      } catch (err) {
          message.error(t('home.uploadFailed'));
      } finally {
          setUploading(false);
      }
  };

  const handleCancel = () => {
      setIsModalOpen(false);
      setFileList([]);
      form.resetFields();
  };

  return (
    <MainLayout {...props}>
      <div style={{ textAlign: 'center', padding: '50px 0' }}>
        <h1>{t('common.home')}</h1>
        <p style={{ fontSize: '18px', color: '#666' }}>Welcome to HMiMG - HanceyMica Image Management Gallery</p>
        
        <div style={{ maxWidth: '600px', margin: '40px auto' }}>
             <Dragger 
                name="file" 
                multiple={true} 
                beforeUpload={handleBeforeUpload}
                showUploadList={false}
                fileList={[]}
             >
                <p className="ant-upload-drag-icon">
                  <InboxOutlined />
                </p>
                <p className="ant-upload-text">{t('home.dragDropText')}</p>
                <p className="ant-upload-hint">
                  {t('home.dragDropHint')}
                </p>
              </Dragger>
        </div>

        <div style={{ marginTop: '40px', display: 'flex', justifyContent: 'center', gap: '20px' }}>
             <Link href="/library">
                <Button type="primary" size="large" icon={<PictureOutlined />}>
                    {t('common.library')}
                </Button>
             </Link>
             <Link href="/admin">
                <Button size="large" icon={<FolderOutlined />}>
                    {t('common.admin')}
                </Button>
             </Link>
        </div>
      </div>

      <Modal
          title={t('home.selectAlbumToUpload')}
          open={isModalOpen}
          onCancel={handleCancel}
          footer={null}
      >
          <p>{fileList.length} files selected</p>
          <Form form={form} onFinish={handleUpload} layout="vertical">
              <Form.Item 
                name="albumId" 
                label={t('admin.album')} 
                rules={[{ required: true, message: t('login.required') + t('admin.album') }]}
              >
                  <Select placeholder={t('home.selectAlbum')} showSearch optionFilterProp="children">
                      {albums.map(a => <Select.Option key={a.id} value={a.id}>{a.name}</Select.Option>)}
                  </Select>
              </Form.Item>
              <Button type="primary" htmlType="submit" loading={uploading} block>
                  {t('home.upload')}
              </Button>
          </Form>
      </Modal>
    </MainLayout>
  );
};

export default Home;
