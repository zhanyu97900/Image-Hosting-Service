// 定义SHA-256加密函数
function sha256Encrypt(password) {
    return CryptoJS.SHA256(password).toString();
  }
  
  // 这里将函数暴露出去，方便外部文件引用
  if (typeof module === 'object' && typeof module.exports === 'object') {
    module.exports = sha256Encrypt;
  } else if (typeof window === 'object') {
    window.sha256Encrypt = sha256Encrypt;
  }