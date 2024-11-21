// // frontend/pages/upload.js
import { useState } from 'react';
import { 
    convertTextToVoice, 
    convertVoiceToText, 
    removeBackground, 
    recognizeSpeech, 
    recognizeFace, 
    performOCR, 
    translateText 
} from '../utils/api';
import axios from 'axios';

export default function Upload() {
    const [text, setText] = useState('');
    const [audioFile, setAudioFile] = useState(null);
    const [imageFile, setImageFile] = useState(null);
    const [destLang, setDestLang] = useState('');
    const [result, setResult] = useState(null);
    const [language, setLanguage] = useState("en");

    const handleTTS = async () => {
        try {
            const response = await convertTextToVoice({ text, language });
            setResult(response.data);
        } catch (error) {
            console.error('Error in Text-to-Voice:', error);
        }
    };

    const handleVTS = async () => {
        if (!audioFile) {
            alert('Please upload an audio file first.');
            return;
        }

        try {
            // Tải file audio lên endpoint /upload-audio để lấy URL
            const formData = new FormData();
            formData.append('audio', audioFile);

            const uploadResponse = await axios.post('http://localhost:81/upload-audio', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            });

            const audioUrl = uploadResponse.data.audio_url;

            const response = await convertVoiceToText(audioUrl);
            setResult(response.data);
        } catch (error) {
            console.error('Error in Voice-to-Text:', error);
        }
    };

    const handleRemoveBg = async () => {
        if (!imageFile) {
            alert('Please upload an image file first.');
            return;
        }

        try {
            const response = await removeBackground(imageFile);
            setResult(response.data);
        } catch (error) {
            console.error('Error in Background Removal:', error);
        }
    };

    const handleSpeechRecognition = async () => {
        if (!audioFile) {
            alert('Please upload an audio file first.');
            return;
        }

        try {
            // Tải file audio lên endpoint /upload-audio để lấy URL
            const formData = new FormData();
            formData.append('audio', audioFile);

            const uploadResponse = await axios.post('http://localhost:81/upload-audio', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            });

            const audioUrl = uploadResponse.data.audio_url;

            const response = await recognizeSpeech(audioUrl);
            setResult(response.data);
        } catch (error) {
            console.error('Error in Speech Recognition:', error);
        }
    };

    const handleFaceRecognition = async () => {
        if (!imageFile) {
            alert('Please upload an image file first.');
            return;
        }

        try {
            const response = await recognizeFace(imageFile);
            setResult(response.data);
        } catch (error) {
            console.error('Error in Face Recognition:', error);
        }
    };

    const handleOCR = async () => {
        if (!imageFile) {
            alert('Please upload an image file first.');
            return;
        }

        try {
            const response = await performOCR(imageFile);
            setResult(response.data);
        } catch (error) {
            console.error('Error in OCR:', error);
        }
    };

    const handleTranslate = async () => {
        try {
            const response = await translateText(text, destLang);
            setResult(response.data);
        } catch (error) {
            console.error('Error in Translation:', error);
        }
    };

    return (
        <div>
            <h1>Upload & Process</h1>
            <div>
                <label htmlFor="language">Language:</label>
                <select
                    id="language"
                    value={language}
                    onChange={(e) => setLanguage(e.target.value)}
                >
                    <option value="en">English</option>
                    <option value="vi">Vietnamese</option>
                    <option value="fr">French</option>
                    {/* Thêm các ngôn ngữ khác tại đây */}
                </select>
            </div>
            <div>
                <h2>Text-to-Voice</h2>
                <input
                    type="text"
                    value={text}
                    onChange={(e) => setText(e.target.value)}
                    placeholder="Enter text"
                />
                <button onClick={handleTTS}>Convert to Voice</button>
            </div>
            <div>
                <h2>Voice-to-Text</h2>
                <input
                    type="file"
                    accept="audio/*"
                    onChange={(e) => setAudioFile(e.target.files[0])}
                />
                <button onClick={handleVTS}>Convert to Text</button>
            </div>
            <div>
                <h2>Background Removal</h2>
                <input
                    type="file"
                    accept="image/*"
                    onChange={(e) => setImageFile(e.target.files[0])}
                />
                <button onClick={handleRemoveBg}>Remove Background</button>
            </div>
            <div>
                <h2>Speech Recognition</h2>
                <input
                    type="file"
                    accept="audio/*"
                    onChange={(e) => setAudioFile(e.target.files[0])}
                />
                <button onClick={handleSpeechRecognition}>Recognize Speech</button>
            </div>
            <div>
                <h2>Face Recognition</h2>
                <input
                    type="file"
                    accept="image/*"
                    onChange={(e) => setImageFile(e.target.files[0])}
                />
                <button onClick={handleFaceRecognition}>Recognize Face</button>
            </div>
            <div>
                <h2>OCR</h2>
                <input
                    type="file"
                    accept="image/*"
                    onChange={(e) => setImageFile(e.target.files[0])}
                />
                <button onClick={handleOCR}>Perform OCR</button>
            </div>
            <div>
                <h2>Translation</h2>
                <input
                    type="text"
                    value={text}
                    onChange={(e) => setText(e.target.value)}
                    placeholder="Enter text"
                />
                <input
                    type="text"
                    value={destLang}
                    onChange={(e) => setDestLang(e.target.value)}
                    placeholder="Destination Language (e.g., 'vi', 'en')"
                />
                <button onClick={handleTranslate}>Translate</button>
            </div>
            {result && (
                <div>
                    <h3>Result:</h3>
                    <pre>{JSON.stringify(result, null, 2)}</pre>
                </div>
            )}
        </div>
    );
}
