const path = require("path");

module.exports = {
  entry: "./src/browser.ts",
  mode: "production",
  stats: {
    errorDetails: true,
  },
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: "ts-loader",
        exclude: /node_modules/,
      },
    ],
  },
  resolve: {
    extensions: [".tsx", ".ts", ".js"],
  },
  output: {
    filename: "fetch-event-source.min.js",
    path: path.resolve(__dirname, "dist"),
    library: "FetchEventSource",
    libraryTarget: "umd",
    globalObject: "this",
  },
};
