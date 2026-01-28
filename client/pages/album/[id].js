import React, { useEffect, useState } from 'react';
import MainLayout from '../../components/MainLayout';
import { Card, List, Breadcrumb, Empty, Image, Button, Modal, Form, Input, message } from 'antd';
import { HomeOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import api from '../../lib/api';
import { useRouter } from 'next/router';
import Cookies from 'js-cookie';
import { useI18n } from '../../lib/i18n';
import Link from 'next/link';

const AlbumDetail = (props) => {
  const router = useRouter();
  const { id } = router.query;
  const [album, setAlbum] = useState(null);
  const [images, setImages] = useState([]);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const { t } = useI18n();
  const [form] = Form.useForm();

  useEffect(() => {
    if (id) {
        fetchData(id);
    }
  }, [id]);

  const fetchData = async (albumId) => {
    try {
      const [albumRes, imagesRes] = await Promise.all([
        api.get(`/albums/${albumId}`),
        api.get(`/images?albumId=${albumId}`)
      ]);
      setAlbum(albumRes.data);
      setImages(imagesRes.data);
      form.setFieldsValue(albumRes.data);
    } catch (err) {
      console.error(err);
    }
  };

  const handleUpdate = async (values) => {
      try {
          await api.put(`/albums/${id}`, values);
          message.success(t('image.updateSuccess'));
          setAlbum({ ...album, ...values });
          setIsEditModalOpen(false);
      } catch (err) {
          message.error(t('image.updateFailed'));
      }
  };

  const handleDelete = () => {
      Modal.confirm({
          title: t('image.confirmDelete'),
          onOk: async () => {
              try {
                  await api.delete(`/albums/${id}`);
                  message.success(t('image.deleteSuccess'));
                  router.push('/?tab=albums');
              } catch (err) {
                  message.error(t('image.deleteFailed'));
              }
          }
      });
  };

  if (!album) return <MainLayout {...props}><div>Loading...</div></MainLayout>;

  return (
    <MainLayout {...props}>
      <Breadcrumb style={{ marginBottom: 16 }}>
        <Breadcrumb.Item href="/">
          <HomeOutlined />
        </Breadcrumb.Item>
        <Breadcrumb.Item href="/?tab=albums">{t('common.albums')}</Breadcrumb.Item>
        <Breadcrumb.Item>{album.name}</Breadcrumb.Item>
      </Breadcrumb>

      <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <h1>{album.name}</h1>
            <p>{album.description}</p>
          </div>
          <div style={{ display: 'flex', gap: 10 }}>
              <Button icon={<EditOutlined />} onClick={() => setIsEditModalOpen(true)}>{t('image.edit')}</Button>
              <Button icon={<DeleteOutlined />} danger onClick={handleDelete}>{t('image.delete')}</Button>
          </div>
      </div>

      {images.length === 0 ? (
          <Empty description="No images found" />
      ) : (
            <List
                grid={{ gutter: 16, xs: 2, sm: 3, md: 4, lg: 6, xl: 6, xxl: 8 }}
                dataSource={images}
                renderItem={(item) => (
                <List.Item>
                    <Link href={`/image/${item.id}`}>
                        <Card
                            hoverable
                            cover={
                                <Image
                                    preview={false}
                                    alt={item.original_name}
                                    src={`http://localhost:3001/${item.path}`}
                                    style={{ height: 200, objectFit: 'cover' }}
                                />
                            }
                        >
                            <Card.Meta title={item.original_name} />
                        </Card>
                    </Link>
                </List.Item>
                )}
            />
      )}
      <Modal
          title={t('image.edit')}
          open={isEditModalOpen}
          onCancel={() => setIsEditModalOpen(false)}
          footer={null}
      >
          <Form form={form} onFinish={handleUpdate} layout="vertical">
              <Form.Item name="name" label={t('admin.name')} rules={[{ required: true }]}>
                  <Input />
              </Form.Item>
              <Form.Item name="description" label={t('admin.description')}>
                  <Input.TextArea />
              </Form.Item>
              <Button type="primary" htmlType="submit" block>{t('image.update')}</Button>
          </Form>
      </Modal>
    </MainLayout>
  );
};

export default AlbumDetail;
