import React, { useEffect, useState } from 'react';
import MainLayout from '../../components/MainLayout';
import { Card, Breadcrumb, Row, Col, Descriptions, Button, Image, Modal, Radio, Input, message } from 'antd';
import { HomeOutlined, DownloadOutlined, ArrowLeftOutlined, LinkOutlined, CopyOutlined } from '@ant-design/icons';
import api from '../../lib/api';
import { useRouter } from 'next/router';
import { useI18n } from '../../lib/i18n';

const ImageDetail = (props) => {
  const router = useRouter();
  const { id } = router.query;
  const [image, setImage] = useState(null);
  const [isLinkModalOpen, setIsLinkModalOpen] = useState(false);
  const [linkType, setLinkType] = useState('direct');
  const { t } = useI18n();

  useEffect(() => {
    if (id) {
        fetchData(id);
    }
  }, [id]);

  const fetchData = async (imageId) => {
    try {
      const res = await api.get(`/images/${imageId}`);
      setImage(res.data);
    } catch (err) {
      console.error(err);
    }
  };

  if (!image) return <MainLayout {...props}><div>Loading...</div></MainLayout>;

  const formatBytes = (bytes, decimals = 2) => {
    if (!+bytes) return '0 Bytes';
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
  };

  const getLink = () => {
      const url = `http://localhost:3001/${image.path}`;
      switch (linkType) {
          case 'markdown': return `![${image.original_name}](${url})`;
          case 'html': return `<img src="${url}" alt="${image.original_name}" />`;
          case 'bbcode': return `[img]${url}[/img]`;
          default: return url;
      }
  };

  const handleCopy = () => {
      navigator.clipboard.writeText(getLink());
      message.success(t('image.copied'));
  };

  return (
    <MainLayout {...props}>
      <Breadcrumb style={{ marginBottom: 16 }}>
        <Breadcrumb.Item href="/">
          <HomeOutlined />
        </Breadcrumb.Item>
        <Breadcrumb.Item href="/?tab=photos">{t('common.photos')}</Breadcrumb.Item>
        <Breadcrumb.Item>{image.original_name}</Breadcrumb.Item>
      </Breadcrumb>

      <div style={{ marginBottom: 16 }}>
          <Button icon={<ArrowLeftOutlined />} onClick={() => router.back()}>
              {t('image.back')}
          </Button>
      </div>

      <Row gutter={24}>
          <Col xs={24} md={16}>
              <Card>
                  <div style={{ display: 'flex', justifyContent: 'center', background: '#f0f2f5', padding: '20px' }}>
                    <Image
                        src={`http://localhost:3001/${image.path}`}
                        alt={image.original_name}
                        style={{ maxHeight: '600px', maxWidth: '100%', objectFit: 'contain' }}
                    />
                  </div>
              </Card>
          </Col>
          <Col xs={24} md={8}>
              <Card title={t('image.details')}>
                  <Descriptions column={1} bordered>
                      <Descriptions.Item label={t('image.originalName')}>{image.original_name}</Descriptions.Item>
                      <Descriptions.Item label={t('image.filename')}>{image.filename}</Descriptions.Item>
                      <Descriptions.Item label={t('image.size')}>{formatBytes(image.size)}</Descriptions.Item>
                      <Descriptions.Item label={t('image.type')}>{image.mimetype}</Descriptions.Item>
                      <Descriptions.Item label={t('image.album')}>{image.album_name || '-'}</Descriptions.Item>
                      <Descriptions.Item label={t('image.uploadedAt')}>{new Date(image.created_at).toLocaleString()}</Descriptions.Item>
                  </Descriptions>
                  <div style={{ marginTop: 20, display: 'flex', flexDirection: 'column', gap: 10 }}>
                      <Button 
                        type="primary" 
                        icon={<DownloadOutlined />} 
                        block 
                        href={`http://localhost:3001/${image.path}`} 
                        download={image.original_name}
                        target="_blank"
                      >
                          {t('image.download')}
                      </Button>
                      <Button 
                        icon={<LinkOutlined />} 
                        block 
                        onClick={() => setIsLinkModalOpen(true)}
                      >
                          {t('image.generateLink')}
                      </Button>
                  </div>
              </Card>
          </Col>
      </Row>

      <Modal
          title={t('image.linkSettings')}
          open={isLinkModalOpen}
          onCancel={() => setIsLinkModalOpen(false)}
          footer={null}
      >
          <div style={{ marginBottom: 16 }}>
              <p>{t('image.linkType')}:</p>
              <Radio.Group value={linkType} onChange={e => setLinkType(e.target.value)}>
                  <Radio.Button value="direct">{t('image.direct')}</Radio.Button>
                  <Radio.Button value="markdown">{t('image.markdown')}</Radio.Button>
                  <Radio.Button value="html">{t('image.html')}</Radio.Button>
                  <Radio.Button value="bbcode">{t('image.bbcode')}</Radio.Button>
              </Radio.Group>
          </div>
          <div style={{ display: 'flex', gap: 8 }}>
              <Input value={getLink()} readOnly />
              <Button type="primary" icon={<CopyOutlined />} onClick={handleCopy}>
                  {t('image.copy')}
              </Button>
          </div>
      </Modal>
    </MainLayout>
  );
};

export default ImageDetail;
