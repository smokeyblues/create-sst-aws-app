const config = {

    STRIPE_KEY: "pk_test_51H5flFItMT7HUULIJU25K7QSB50g8qqrwLVmK6XST0IrYivMZnU9fL90LBPkmJl4ptLbtlrC0Ux6fFN9K4xmfu8200oZDQbTJl",

    // Frontend config
    MAX_ATTACHMENT_SIZE: 5000000,

    // Backend config
    s3: {
      REGION: import.meta.env.VITE_REGION,
      BUCKET: import.meta.env.VITE_BUCKET,
    },
    apiGateway: {
      REGION: import.meta.env.VITE_REGION,
      URL: import.meta.env.VITE_API_URL,
    },
    cognito: {
      REGION: import.meta.env.VITE_REGION,
      USER_POOL_ID: import.meta.env.VITE_USER_POOL_ID,
      APP_CLIENT_ID: import.meta.env.VITE_USER_POOL_CLIENT_ID,
      IDENTITY_POOL_ID: import.meta.env.VITE_IDENTITY_POOL_ID,
    },
  };
  
  export default config;