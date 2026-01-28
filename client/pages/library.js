import React, { useEffect, useState } from 'react';
import MainLayout from '../components/MainLayout';
import { Tabs, Button, Modal, Form, Input, Select, Upload, message, Card, List, Empty } from 'antd';
import { UploadOutlined, FolderOutlined, PictureOutlined } from '@ant-design/icons';
import api from '../lib/api';
import { useRouter } from 'next/router';
import Cookies from 'js-cookie';
import { useI18n } from '../lib/i18n';

import Link from 'next/link';

const { Meta } = Card;

const Library = (props) => {
  const router = useRouter();
  const [albums, setAlbums] = useState([]);
  const [collections, setCollections] = useState([]);
  const [allImages, setAllImages] = useState([]);
  const [isUploadModalOpen, setIsUploadModalOpen] = useState(false);
  const [activeTab, setActiveTab] = useState('1');
  const { t } = useI18n();

  useEffect(() => {
    const token = Cookies.get('token');
    if (!token) {
      router.push('/login');
      return;
    }
    const { tab } = router.query;
    if (tab === 'albums') setActiveTab('1');
    else if (tab === 'collections') setActiveTab('2');
    else if (tab === 'photos') setActiveTab('3');
    
    fetchData();
  }, [router.query]);

  const fetchData = async () => {
    try {
      const [albumsRes, collectionsRes, imagesRes] = await Promise.all([
        api.get('/albums'),
        api.get('/collections'),
        api.get('/images')
      ]);
      setAlbums(albumsRes.data);
      setCollections(collectionsRes.data);
      setAllImages(imagesRes.data);
    } catch (err) {
      console.error(err);
    }
  };

  const handleUpload = async (values) => {
    const formData = new FormData();
    if (values.image && values.image.length > 0) {
        values.image.forEach(file => {
            formData.append('images', file.originFileObj);
        });
    }
    formData.append('albumId', values.albumId);

    try {
      await api.post('/upload', formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      });
      message.success(t('home.uploadSuccess'));
      setIsUploadModalOpen(false);
    } catch (err) {
      message.error(t('home.uploadFailed'));
    }
  };

  const items = [
    {
      key: '1',
      label: t('common.albums'),
      children: (
        <List
          grid={{ gutter: 16, xs: 1, sm: 2, md: 3, lg: 4, xl: 4, xxl: 6 }}
          dataSource={albums}
          renderItem={(item) => (
            <List.Item>
              <Link href={`/album/${item.id}`}>
                <Card
                    hoverable
                    cover={
                        item.cover_image ? 
                        <img alt={item.name} src={`http://localhost:3001/${item.cover_image}`} style={{height: 150, objectFit: 'cover'}} /> :
                        <div style={{height: 150, background: '#eee', display: 'flex', alignItems: 'center', justifyContent: 'center'}}><PictureOutlined style={{fontSize: 30}}/></div>
                    }
                >
                    <Meta title={item.name} description={item.description} />
                </Card>
              </Link>
            </List.Item>
          )}
        />
      ),
    },
    {
      key: '2',
      label: t('common.collections'),
      children: (
        <List
          grid={{ gutter: 16, xs: 1, sm: 2, md: 3, lg: 4, xl: 4, xxl: 6 }}
          dataSource={collections}
          renderItem={(item) => (
            <List.Item>
              <Link href={`/collection/${item.id}`}>
                <Card hoverable>
                    <Meta 
                        avatar={<FolderOutlined style={{fontSize: 24}} />}
                        title={item.name} 
                        description={item.description} 
                    />
                </Card>
              </Link>
            </List.Item>
          )}
        />
      ),
    },
    {
      key: '3',
      label: t('common.photos'),
      children: (
        <List
          grid={{ gutter: 16, xs: 2, sm: 3, md: 4, lg: 6, xl: 6, xxl: 8 }}
          dataSource={allImages}
          renderItem={(item) => (
            <List.Item>
              <Link href={`/image/${item.id}`}>
                <Card
                    hoverable
                    cover={
                        <img alt={item.original_name} src={`http://localhost:3001/${item.path}`} style={{height: 150, objectFit: 'cover'}} />
                    }
                >
                    <Meta title={item.original_name} />
                </Card>
              </Link>
            </List.Item>
          )}
        />
      ),
    },
  ];

  return (
    <MainLayout {...props}>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2>{t('home.myLibrary')}</h2>
        <Button type="primary" icon={<UploadOutlined />} onClick={() => setIsUploadModalOpen(true)}>
          {t('home.uploadPhoto')}
        </Button>
      </div>

      <Tabs defaultActiveKey="1" items={items} onChange={setActiveTab} />

      <Modal
        title={t('home.uploadPhoto')}
        open={isUploadModalOpen}
        onCancel={() => setIsUploadModalOpen(false)}
        footer={null}
      >
        <Form onFinish={handleUpload} layout="vertical">
          <Form.Item
            name="albumId"
            label={t('home.selectAlbum')}
            rules={[{ required: true, message: t('login.required') + t('admin.album') }]}
          >
            <Select placeholder={t('home.selectAlbum')}>
              {albums.map(album => (
                <Select.Option key={album.id} value={album.id}>{album.name}</Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="image"
            label={t('home.selectFile')}
            valuePropName="fileList"
            getValueFromEvent={(e) => {
              if (Array.isArray(e)) return e;
              return e?.fileList;
            }}
            rules={[{ required: true, message: t('login.required') + t('home.selectFile') }]}
          >
            <Upload beforeUpload={() => false} multiple listType="picture">
              <Button icon={<UploadOutlined />}>{t('home.selectFile')}</Button>
            </Upload>
          </Form.Item>
          <Button type="primary" htmlType="submit" block>
            {t('home.upload')}
          </Button>
        </Form>
      </Modal>
    </MainLayout>
  );
};

export default Library;
