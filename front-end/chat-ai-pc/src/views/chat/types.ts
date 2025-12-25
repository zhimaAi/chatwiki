export interface TextMessage {
  type: "text";
  uid: string;
  text: string;
}

export interface ImageUrlMessage {
  uid: string;
  type: 'image_url';
  image_url: {
    url: string;
  };
}

export type Message = TextMessage | ImageUrlMessage;