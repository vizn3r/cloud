function showNotification(message, type = 'success') {
  const notification = document.getElementById('notification');
  const messageEl = document.getElementById('notificationMessage');

  messageEl.textContent = message;
  notification.className = `notification ${type} show`;

  // Auto-hide after 3 seconds
  setTimeout(() => {
    notification.className = 'notification';
  }, 3000);
}

function formatFileSize(bytes) {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i];
}

function toggleFileDetails(fileId) {
  const detailsElement = document.getElementById(`details-${fileId}`);
  const button = detailsElement.previousElementSibling.querySelector('.btn-info');

  if (detailsElement.style.display === 'none') {
    detailsElement.style.display = 'block';
    button.textContent = 'Show less';
  } else {
    detailsElement.style.display = 'none';
    button.textContent = 'Show more';
  }
}

// Check authentication and load files on page load
document.addEventListener('DOMContentLoaded', function() {
  const token = localStorage.getItem('authToken');

  if (!token) {
    window.location.href = '/login';
    return;
  }

  verifyTokenAndLoadFiles();
});

function verifyTokenAndLoadFiles() {
  const token = localStorage.getItem('authToken');

  // Verify token and get user info
  fetch('http://localhost:8080/user/me', {
    headers: {
      'Authorization': 'Bearer ' + token
    }
  })
  .then(response => {
    if (!response.ok) {
      localStorage.removeItem('authToken');
      localStorage.removeItem('userId');
      window.location.href = '/login';
      return;
    }
    return response.json();
  })
  .then(userInfo => {
    if (userInfo) {
      document.getElementById('userEmail').textContent = userInfo.email;
      loadUserFiles();
    }
  })
  .catch(error => {
    console.error('Auth check failed:', error);
    localStorage.removeItem('authToken');
    localStorage.removeItem('userId');
    window.location.href = '/login';
  });
}

async function loadUserFiles() {
  const token = localStorage.getItem('authToken');

  try {
    // First get the list of file IDs
    const filesResponse = await fetch('http://localhost:8080/user/files', {
      headers: {
        'Authorization': 'Bearer ' + token
      }
    });

    if (!filesResponse.ok) {
      throw new Error('Failed to load files');
    }

    const data = await filesResponse.json();
    const fileIds = data.files || [];

    // Fetch metadata for each file
    userFiles = await Promise.all(fileIds.map(async (fileId) => {
      try {
        const metaResponse = await fetch(`http://localhost:8080/file/${fileId}/data`, {
          headers: {
            'Authorization': 'Bearer ' + token
          }
        });

        if (!metaResponse.ok) {
          throw new Error(`Failed to load metadata for ${fileId}`);
        }

        const metadata = await metaResponse.json();
        return { id: fileId, ...metadata };
      } catch (error) {
        console.error(`Error loading metadata for ${fileId}:`, error);
        return { id: fileId, uploadName: fileId, size: 0, contentType: 'unknown' };
      }
    }));

    renderFilesGallery().catch(error => {
      console.error('Error rendering files gallery:', error);
    });
  } catch (error) {
    console.error('Error loading files:', error);
    document.getElementById('filesContainer').innerHTML =
      '<div class="error">Error loading files: ' + error.message + '</div>';
  }
}

async function loadThumbnailWithAuth(fileId) {
  const token = localStorage.getItem('authToken');
  if (!token) {
    throw new Error('No authentication token');
  }

  try {
    const response = await fetch(`http://localhost:8080/file/${fileId}/thumbnail`, {
      headers: {
        'Authorization': 'Bearer ' + token
      }
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const blob = await response.blob();
    return URL.createObjectURL(blob);
  } catch (error) {
    console.error(`Failed to load thumbnail for ${fileId}:`, error);
    throw error;
  }
}

async function renderFilesGallery() {
  const container = document.getElementById('filesContainer');

  if (userFiles.length === 0) {
    container.innerHTML = '<div class="empty-state">No files uploaded yet</div>';
    return;
  }

  // Create initial HTML with loading state
  container.innerHTML = userFiles.map(file => `
    <div class="file-card" data-file-id="${file.id}">
      <div class="file-thumbnail" onclick="previewFile('${file.id}')">
        <div class="thumbnail-loading" style="width: 100%; height: 250px; display; flex; align-items: center; justify-content: center; background: #f3f4f6;">
          <span style="font-size: 18px; color: #6b7280;">Loading...</span>
        </div>
        <img class="thumbnail-image" style="display: none;" alt="File thumbnail">
        <div class="file-icon" style="display: none;">📄</div>
      </div>
      <div class="file-info">
        <div class="file-name">${file.uploadName || file.id}</div>
        <div class="file-actions">
          <button onclick="downloadFile('${file.id}')" class="btn-action btn-download">Download</button>
          <button onclick="shareFile('${file.id}')" class="btn-action btn-share">Share</button>
          <button onclick="deleteFile('${file.id}')" class="btn-action btn-delete">Delete</button>
          <button onclick="toggleFileDetails('${file.id}')" class="btn-action btn-info">Show more</button>
        </div>
        <div class="file-details" id="details-${file.id}" style="display: none;">
          <div class="detail-row"><strong>ID:</strong> ${file.id}</div>
          <div class="detail-row"><strong>Size:</strong> ${formatFileSize(file.size)}</div>
          <div class="detail-row"><strong>Type:</strong> ${file.contentType || 'unknown'}</div>
          <div class="detail-row"><strong>Uploaded:</strong> ${file.uploadedAt || 'Unknown date'}</div>
        </div>
      </div>
    </div>
  `).join('');

  // Load thumbnails asynchronously
  userFiles.forEach(async (file, index) => {
    const card = container.children[index];
    const thumbnailImg = card.querySelector('.thumbnail-image');
    const loadingEl = card.querySelector('.thumbnail-loading');
    const fileIcon = card.querySelector('.file-icon');

    try {
      const thumbnailUrl = await loadThumbnailWithAuth(file.id);
      thumbnailImg.src = thumbnailUrl;
      thumbnailImg.style.display = 'block';
      loadingEl.style.display = 'none';
    } catch (error) {
      // If thumbnail loading fails, show file icon instead
      loadingEl.style.display = 'none';
      fileIcon.style.display = 'block';
    }
  });
}

function showUploadModal() {
  const modal = document.getElementById('uploadModal');
  const fileInput = document.getElementById('fileInput');
  const modalFileInput = document.getElementById('modalFileInput');

  // Transfer files from floating button to modal
  if (fileInput.files.length > 0) {
    // Create a new FileList-like object
    const dataTransfer = new DataTransfer();
    Array.from(fileInput.files).forEach(file => {
      dataTransfer.items.add(file);
    });
    modalFileInput.files = dataTransfer.files;
  }

  modal.classList.add('show');
}

function closeUploadModal() {
  const modal = document.getElementById('uploadModal');
  modal.classList.remove('show');
}

function uploadFile() {
  const fileInput = document.getElementById('modalFileInput');
  const files = fileInput.files;

  if (!files || files.length === 0) {
    showNotification('Please select at least one file', 'error');
    return;
  }

  const token = localStorage.getItem('authToken');
  const responseEl = document.getElementById('uploadResponse');

  // Upload each file
  Array.from(files).forEach((file, index) => {
    const formData = new FormData();
    formData.append('file', file);

    fetch('http://localhost:8080/file', {
      method: 'POST',
      headers: {
        'Authorization': 'Bearer ' + token
      },
      body: formData
    })
    .then(response => response.text())
    .then(fileId => {
      responseEl.innerHTML = `File "${file.name}" uploaded successfully! File ID: ${fileId}`;
      responseEl.className = 'response success';

      // Reload files after upload and close modal
      if (index === files.length - 1) {
        loadUserFiles();
        setTimeout(() => {
          closeUploadModal();
        }, 1000);
      }
    })
    .catch(error => {
      responseEl.innerHTML = `Error uploading "${file.name}": ` + error.message;
      responseEl.className = 'response error';
    });
  });
}

function downloadFile(fileId) {
  const token = localStorage.getItem('authToken');

  fetch(`http://localhost:8080/file/${fileId}`, {
    headers: {
      'Authorization': 'Bearer ' + token
    }
  })
  .then(response => {
    if (!response.ok) {
      throw new Error('Failed to download file');
    }
    return response.blob();
  })
  .then(blob => {
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `file-${fileId}`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  })
  .catch(error => {
    showNotification('Error downloading file: ' + error.message, 'error');
  });
}

function shareFile(fileId) {
  currentShareFileId = fileId;
  const token = localStorage.getItem('authToken');

  document.getElementById('shareFileName').textContent = fileId;
  document.getElementById('shareModal').classList.add('show');

  // Create share link
  fetch(`http://localhost:8080/share/${fileId}`, {
    method: 'POST',
    headers: {
      'Authorization': 'Bearer ' + token,
      'X-Share-Duration': '10080' // 1 week in minutes
    }
  })
  .then(response => response.text())
  .then(shareId => {
    const shareLink = `http://localhost:8080/share/${shareId}`;
    document.getElementById('shareLinkInput').value = shareLink;
  })
  .catch(error => {
    showNotification('Error creating share: ' + error.message, 'error');
    closeShareModal();
  });
}

function closeShareModal() {
  document.getElementById('shareModal').classList.remove('show');
  currentShareFileId = null;
}

function copyShareLink() {
  const shareLinkInput = document.getElementById('shareLinkInput');
  shareLinkInput.select();
  document.execCommand('copy');
  showNotification('Share link copied to clipboard!');
}

function deleteFile(fileId) {
  currentDeleteFileId = fileId;
  document.getElementById('deleteModal').classList.add('show');
}

function confirmDelete() {
  if (!currentDeleteFileId) return;

  const token = localStorage.getItem('authToken');
  const fileId = currentDeleteFileId;

  fetch(`http://localhost:8080/file/${fileId}`, {
    method: 'DELETE',
    headers: {
      'Authorization': 'Bearer ' + token
    }
  })
  .then(response => {
    if (response.ok) {
      // Remove from local list and re-render
      userFiles = userFiles.filter(id => id !== fileId);
      renderFilesGallery().catch(error => {
        console.error('Error rendering files gallery:', error);
      });
      showNotification('File deleted successfully');
    } else {
      throw new Error('Failed to delete file');
    }
  })
  .catch(error => {
    showNotification('Error deleting file: ' + error.message, 'error');
  })
  .finally(() => {
    cancelDelete();
  });
}

function cancelDelete() {
  currentDeleteFileId = null;
  document.getElementById('deleteModal').classList.remove('show');
}

function logout() {
  localStorage.removeItem('authToken');
  localStorage.removeItem('userId');
  window.location.href = '/login';
}

let currentPreviewFileId = null;

function previewFile(fileId) {
  currentPreviewFileId = fileId;
  const token = localStorage.getItem('authToken');

  document.getElementById('previewFileName').textContent = fileId;
  document.getElementById('previewContent').innerHTML = '<div style="color: white; text-align: center; padding: 2rem;">Loading preview...</div>';
  document.getElementById('previewOverlay').style.display = 'flex';

  // Check file type and load appropriate preview
  fetch(`http://localhost:8080/file/${fileId}`, {
    headers: {
      'Authorization': 'Bearer ' + token
    }
  })
  .then(response => {
    const contentType = response.headers.get('content-type') || '';

    if (contentType.includes('image')) {
      return response.blob().then(blob => {
        const url = URL.createObjectURL(blob);
        document.getElementById('previewContent').innerHTML = `
          <img src="${url}" style="max-width: 100vw; max-height: 100vh; object-fit: contain;">
        `;
      });
    } else if (contentType.includes('pdf')) {
      return response.blob().then(blob => {
        const url = URL.createObjectURL(blob);
        document.getElementById('previewContent').innerHTML = `
          <embed src="${url}" type="application/pdf" width="100%" height="100%">
        `;
      });
    } else if (contentType.includes('text')) {
      return response.text().then(text => {
        document.getElementById('previewContent').innerHTML = `
          <pre style="white-space: pre-wrap; background: #1f2937; color: #e5e7eb; padding: 2rem; border-radius: 8px; max-width: 90%; max-height: 90%; overflow: auto; font-family: 'Monaco', 'SF Mono', monospace; font-size: 0.9rem; line-height: 1.5;">${text}</pre>
        `;
      });
    } else {
      document.getElementById('previewContent').innerHTML = `
        <div style="color: white; text-align: center; padding: 2rem;">
          <p>No preview available for this file type</p>
          <button onclick="downloadFile('${fileId}')" class="btn btn-download">Download File</button>
        </div>
      `;
    }
  })
  .catch(error => {
    console.error('Preview error:', error);
    document.getElementById('previewContent').innerHTML = `
      <div style="text-align: center; padding: 2rem; color: white;">
        Error loading preview: ${error.message}
      </div>
    `;
  });
}

function closePreviewModal() {
  document.getElementById('previewOverlay').style.display = 'none';
  currentPreviewFileId = null;
}

// Close modals when clicking outside
window.onclick = function(event) {
  const shareModal = document.getElementById('shareModal');
  const deleteModal = document.getElementById('deleteModal');
  const previewModal = document.getElementById('previewModal');

  if (event.target === shareModal) {
    closeShareModal();
  }
  if (event.target === deleteModal) {
    cancelDelete();
  }
  if (event.target === previewModal) {
    closePreviewModal();
  }
}

