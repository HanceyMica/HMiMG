import React, { useEffect, useState } from 'react';
import MainLayout from '../../components/MainLayout';
import { Card, List, Breadcrumb, Empty, Button, Modal, Form, Input, message } from 'antd';
import { HomeOutlined, FolderOutlined, PictureOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import api from '../../lib/api';
import { useRouter } from 'next/router';
import { useI18n } from '../../lib/i18n';
import Link from 'next/link';

const CollectionDetail = (props) => {
  const router = useRouter();
  const { id } = router.query;
  const [collection, setCollection] = useState(null);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const { t } = useI18n();
  const [form] = Form.useForm();

  useEffect(() => {
    if (id) {
        fetchData(id);
    }
  }, [id]);

  const fetchData = async (collectionId) => {
    try {
      const res = await api.get(`/collections/${collectionId}`);
      setCollection(res.data);
      form.setFieldsValue(res.data);
    } catch (err) {
      console.error(err);
    }
  };

  const handleUpdate = async (values) => {
      try {
          await api.put(`/collections/${id}`, values);
          message.success(t('image.updateSuccess'));
          setCollection({ ...collection, ...values });
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
                  await api.delete(`/collections/${id}`);
                  message.success(t('image.deleteSuccess'));
                  router.push('/?tab=collections');
              } catch (err) {
                  message.error(t('image.deleteFailed'));
              }
          }
      });
  };

  if (!collection) return <MainLayout {...props}><div>Loading...</div></MainLayout>;

  return (
    <MainLayout {...props}>
      <Breadcrumb style={{ marginBottom: 16 }}>
        <Breadcrumb.Item href="/">
          <HomeOutlined />
        </Breadcrumb.Item>
        <Breadcrumb.Item href="/?tab=collections">{t('common.collections')}</Breadcrumb.Item>
        <Breadcrumb.Item>{collection.name}</Breadcrumb.Item>
      </Breadcrumb>

      <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <h1>{collection.name}</h1>
            <p>{collection.description}</p>
          </div>
          <div style={{ display: 'flex', gap: 10 }}>
              <Button icon={<EditOutlined />} onClick={() => setIsEditModalOpen(true)}>{t('image.edit')}</Button>
              <Button icon={<DeleteOutlined />} danger onClick={handleDelete}>{t('image.delete')}</Button>
          </div>
      </div>

      {!collection.children || collection.children.length === 0 ? (
          <Empty description="No items found" />
      ) : (
          <List
            grid={{ gutter: 16, xs: 1, sm: 2, md: 3, lg: 4, xl: 4, xxl: 6 }}
            dataSource={collection.children}
            renderItem={(item) => (
                <List.Item>
                    <Link href={item.type === 'album' ? `/album/${item.id}` : `/collection/${item.id}`}>
                        <Card
                            hoverable
                            cover={
                                item.type === 'album' ? (
                                    item.cover_image ? 
                                    <img alt={item.name} src={`http://localhost:3001/${item.cover_image}`} style={{height: 150, objectFit: 'cover'}} /> :
                                    <div style={{height: 150, background: '#eee', display: 'flex', alignItems: 'center', justifyContent: 'center'}}><PictureOutlined style={{fontSize: 30}}/></div>
                                ) : (
                                    <div style={{height: 150, background: '#f0f2f5', display: 'flex', alignItems: 'center', justifyContent: 'center'}}><FolderOutlined style={{fontSize: 40}}/></div>
                                )
                            }
                        >
                            <Card.Meta 
                                avatar={item.type === 'collection' ? <FolderOutlined /> : <PictureOutlined />}
                                title={item.name} 
                                description={item.description} 
                            />
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

export default CollectionDetail;
