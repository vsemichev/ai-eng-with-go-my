module.exports = {
  webpack: {
    configure: (webpackConfig) => {
      // MUI ships dual CJS/ESM builds. When webpack resolves the "module"
      // (ESM) field it sometimes misdetects named exports on files that use
      // dynamic Object.defineProperty-based exports, producing false
      // "does not contain a default export" build failures. Preferring the
      // CJS "main" field avoids that false positive without affecting
      // runtime behavior.
      webpackConfig.resolve.mainFields = ['main', 'module'];
      return webpackConfig;
    },
  },
};
