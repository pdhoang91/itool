// frontend/utils/api.js

import axios from 'axios';

const API_BASE_URL = 'http://localhost:81';

export const convertTextToVoice = (text) => {
    return axios.post(`${API_BASE_URL}/tts`, { text });
};

export const convertVoiceToText = (audioUrl) => {
    return axios.post(`${API_BASE_URL}/vts`, { audio_url: audioUrl });
};

export const removeBackground = (imageFile) => {
    const formData = new FormData();
    formData.append('image', imageFile);
    return axios.post(`${API_BASE_URL}/remove-bg`, formData, {
        headers: {
            'Content-Type': 'multipart/form-data',
        },
    });
};

export const recognizeSpeech = (audioUrl) => {
    return axios.post(`${API_BASE_URL}/speech-recognition`, { audio_url: audioUrl });
};

export const recognizeFace = (imageFile) => {
    const formData = new FormData();
    formData.append('image', imageFile);
    return axios.post(`${API_BASE_URL}/face-recognition`, formData, {
        headers: {
            'Content-Type': 'multipart/form-data',
        },
    });
};

export const performOCR = (imageFile) => {
    const formData = new FormData();
    formData.append('image', imageFile);
    return axios.post(`${API_BASE_URL}/ocr`, formData, {
        headers: {
            'Content-Type': 'multipart/form-data',
        },
    });
};

export const translateText = (text, destLang) => {
    return axios.post(`${API_BASE_URL}/translate`, { text, dest_lang: destLang });
};

export const getTaskStatus = (taskId) => {
    return axios.get(`${API_BASE_URL}/tasks/${taskId}`);
};
