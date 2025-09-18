class FileEncryption {
  static async generateCEK() {
    return await crypto.subtle.generateKey(
      { name: "AES-GCM", length: 256 },
      true, // extractable
      ["encrypt", "decrypt"]
    );
  }

  static async deriveUserKey(password, salt) {
    const keyMaterial = await crypto.subtle.importKey(
      "raw",
      new TextEncoder().encode(password),
      "PBKDF2",
      false,
      ["deriveKey"]
    );

    return await crypto.subtle.deriveKey(
      {
        name: "PBKDF2",
        salt: salt,
        iterations: 100000,
        hash: "SHA-256"
      },
      keyMaterial,
      { name: "AES-GCM", length: 256 },
      false,
      ["wrapKey", "unwrapKey"]
    );
  }

  static async encryptFile(fileData, cek) {
    const iv = crypto.getRandomValues(new Uint8Array(12));
    const encrypted = await crypto.subtle.encrypt(
      { name: "AES-GCM", iv },
      cek,
      fileData
    );
    return { encrypted, iv };
  }

  static async decryptFile(encryptedData, cek, iv) {
    return await crypto.subtle.decrypt(
      { name: "AES-GCM", iv },
      cek,
      encryptedData
    );
  }

  static async wrapCEK(cek, kek) {
    const iv = crypto.getRandomValues(new Uint8Array(12));
    const cekRaw = await crypto.subtle.exportKey("raw", cek);
    const wrapped = await crypto.subtle.encrypt(
      { name: "AES-GCM", iv },
      kek,
      cekRaw
    );
    return { wrapped, iv };
  }

  static async unwrapCEK(wrappedCEK, kek, iv) {
    const cekRaw = await crypto.subtle.decrypt(
      { name: "AES-GCM", iv },
      kek,
      wrappedCEK
    );
    return await crypto.subtle.importKey(
      "raw",
      cekRaw,
      { name: "AES-GCM", length: 256 },
      true,
      ["encrypt", "decrypt"]
    );
  }

  static generateRandomBytes(length) {
    return crypto.getRandomValues(new Uint8Array(length));
  }

  static async arrayBufferToBase64(buffer) {
    const bytes = new Uint8Array(buffer);
    let binary = '';
    for (let i = 0; i < bytes.byteLength; i++) {
      binary += String.fromCharCode(bytes[i]);
    }
    return btoa(binary);
  }

  static async base64ToArrayBuffer(base64) {
    const binary = atob(base64);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
      bytes[i] = binary.charCodeAt(i);
    }
    return bytes;
  }
}
