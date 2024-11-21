// frontend/services/api.js

import axios from 'axios';

const API_BASE_URL = 'http://localhost:81';

// Cấu hình axios instance
const axiosInstance = axios.create({
    baseURL: API_BASE_URL,
});

// Thêm interceptor để thêm token vào headers nếu cần
axiosInstance.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

// API Calls
export const fetchTasks = () => axiosInstance.get('/tasks');

export const convertTextToVoice = ({ text, language }) =>
    axiosInstance.post('/tts', { text, language });

export const convertVoiceToText = (audioUrl) =>
    axiosInstance.post('/vts', { audio_url: audioUrl });

export const removeBackground = (imageFile) => {
    const formData = new FormData();
    formData.append('image', imageFile);
    return axiosInstance.post('/remove-bg', formData, {
        headers: {
            'Content-Type': 'multipart/form-data',
        },
    });
};

export const recognizeSpeech = (audioUrl) =>
    axiosInstance.post('/speech-recognition', { audio_url: audioUrl });

export const recognizeFace = (imageFile) => {
    const formData = new FormData();
    formData.append('image', imageFile);
    return axiosInstance.post('/face-recognition', formData, {
        headers: {
            'Content-Type': 'multipart/form-data',
        },
    });
};

export const performOCR = (imageFile) => {
    const formData = new FormData();
    formData.append('image', imageFile);
    return axiosInstance.post('/ocr', formData, {
        headers: {
            'Content-Type': 'multipart/form-data',
        },
    });
};

export const translateText = (text, destLang) =>
    axiosInstance.post('/translate', { text, dest_lang: destLang });

export const uploadAudio = (audioFile) => {
    const formData = new FormData();
    formData.append('audio', audioFile);
    return axiosInstance.post('/upload-audio', formData, {
        headers: {
            'Content-Type': 'multipart/form-data',
        },
    });
};
