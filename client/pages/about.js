import React from 'react';
import MainLayout from '../components/MainLayout';
import { Card, Typography, List, Divider, Space, Tag, Button, Avatar } from 'antd';
import { GithubOutlined, UserOutlined, CheckCircleFilled, RocketOutlined } from '@ant-design/icons';
import { useI18n } from '../lib/i18n';
import pkg from '../package.json';

const { Title, Paragraph, Text } = Typography;

const About = (props) => {
  const { t } = useI18n();
  const { isDarkMode } = props;

  const features = [
    t('about.feature1'),
    t('about.feature2'),
    t('about.feature3'),
    t('about.feature4'),
  ];

  const sloganImg = isDarkMode ? '/images/slogan_dark.png' : '/images/slogan_light.png';

  return (
    <MainLayout {...props}>
      <div style={{ maxWidth: 800, margin: '0 auto', padding: '40px 0' }}>
        <Card 
            styles={{ body: { border: 'none' } }}
            style={{ 
                textAlign: 'center', 
                overflow: 'hidden',
                boxShadow: '0 8px 32px 0 rgba(31, 38, 135, 0.37)'
            }}
        >
            <div style={{ padding: '0 40px' }}>
                <Title level={2} style={{ marginBottom: 20 }}>{t('about.title')}</Title>
                <Paragraph type="secondary" style={{ fontSize: 18, marginBottom: 40 }}>
                    {t('about.description')}
                </Paragraph>

            {/* Slogan Banner */}
            <div style={{ marginBottom: 40, padding: '20px 0' }}>
                <img 
                    src={sloganImg} 
                    alt="HMiMG Slogan" 
                    style={{ 
                        maxWidth: '80%', 
                        height: 'auto', 
                        maxHeight: 480,
                        objectFit: 'contain'
                    }} 
                />
            </div>

                <div style={{ textAlign: 'left', background: isDarkMode ? 'rgba(255,255,255,0.05)' : 'rgba(0,0,0,0.02)', padding: 30, borderRadius: 16 }}>
                    <Title level={4} style={{ textAlign: 'center', marginBottom: 30 }}>
                        <RocketOutlined /> {t('about.features')}
                    </Title>
                    <List
                        grid={{ gutter: 16, column: 2, xs: 1 }}
                        dataSource={features}
                        renderItem={item => (
                            <List.Item>
                                <Space align="start">
                                    <CheckCircleFilled style={{ color: '#52c41a', fontSize: 18, marginTop: 4 }} />
                                    <Text style={{ fontSize: 16 }}>{item}</Text>
                                </Space>
                            </List.Item>
                        )}
                    />
                </div>

                <Divider style={{ margin: '40px 0' }} />

                <div style={{ display: 'flex', justifyContent: 'center', gap: 40, flexWrap: 'wrap' }}>
                    <div style={{ textAlign: 'center' }}>
                        <Text type="secondary" style={{ display: 'block', marginBottom: 8 }}>{t('about.version')}</Text>
                        <Tag color="blue" style={{ fontSize: 14, padding: '4px 12px' }}>v{pkg.version}</Tag>
                    </div>
                    <div style={{ textAlign: 'center' }}>
                        <Text type="secondary" style={{ display: 'block', marginBottom: 8 }}>{t('about.author')}</Text>
                        <Space>
                            <Avatar icon={<UserOutlined />} style={{ backgroundColor: '#87d068' }} />
                            <Text strong>HanceyMica</Text>
                        </Space>
                    </div>
                    <div style={{ textAlign: 'center' }}>
                        <Text type="secondary" style={{ display: 'block', marginBottom: 8 }}>{t('about.github')}</Text>
                        <Button 
                            type="text" 
                            icon={<GithubOutlined style={{ fontSize: 22 }} />} 
                            href="https://github.com/HanceyMica/HMiMG" 
                            target="_blank"
                        />
                    </div>
                </div>
                
                <div style={{ marginTop: 40, color: '#999' }}>
                     {t('common.copyright')}
                </div>
            </div>
        </Card>
      </div>
    </MainLayout>
  );
};

export default About;
