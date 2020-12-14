import { Router } from 'express';
import webhookRouter from './webhook';

const router = Router();

router.use('/webhook', webhookRouter);

router.get('/healthcheck', (req, res, next) => {
  res.json({});
});

export default router;
