from flask import Flask, request, jsonify
import face_recognition
import os

app = Flask(__name__)

@app.route('/recognize-face', methods=['POST'])
def recognize_face():
    if 'image' not in request.files:
        return jsonify({"error": "No image file provided"}), 400

    image = request.files['image']
    image_path = f"/tmp/{image.filename}"
    image.save(image_path)

    try:
        # Thực hiện nhận diện khuôn mặt
        img = face_recognition.load_image_file(image_path)
        face_locations = face_recognition.face_locations(img)

        response = {"face_count": len(face_locations)}
    except Exception as e:
        response = {"error": str(e)}

    # Xóa file tạm
    os.remove(image_path)

    return jsonify(response), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5005)
