import React from 'react';
import MainLayout from '../components/MainLayout';
import { Result, Button } from 'antd';
import { useRouter } from 'next/router';
import { useI18n } from '../lib/i18n';

const Custom404 = () => {
  const router = useRouter();
  const { t } = useI18n();

  return (
    <MainLayout>
      <Result
        status="404"
        title={t('notFound.title')}
        subTitle={t('notFound.description')}
        extra={
          <Button type="primary" onClick={() => router.push('/')}>
            {t('notFound.backHome')}
          </Button>
        }
      />
    </MainLayout>
  );
};

export default Custom404;
