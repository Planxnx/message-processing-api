import { Router } from 'express';
const router = Router();

router.post('/message-processing-api', async (req, res, next) => {
  console.log(req.body)
  res.json({});
});

export default router;
