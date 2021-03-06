interface WeatherCode {
  code: number;
  text: string;
}

class Weather {
  static weatherCodes: WeatherCode[] = [
    {code: 100, text: "晴れ"},
    {code: 101, text: "晴れ 時々 くもり"},
    {code: 102, text: "晴れ 一時 雨"},
    {code: 103, text: "晴れ 時々 雨"},
    {code: 104, text: "晴れ 一時 雪"},
    {code: 105, text: "晴れ 時々 雪"},
    {code: 106, text: "晴れ 一時 雨か雪"},
    {code: 107, text: "晴れ 時々 雨か雪"},
    {code: 108, text: "晴れ 一時 雨か雷雨"},
    {code: 110, text: "晴れ のち時々くもり"},
    {code: 111, text: "晴れ のち くもり"},
    {code: 112, text: "晴れ のち一時 雨"},
    {code: 113, text: "晴れ のち時々 雨"},
    {code: 114, text: "晴れ のち 雨"},
    {code: 115, text: "晴れ のち一時 雪"},
    {code: 116, text: "晴れ のち時々 雪"},
    {code: 117, text: "晴れ のち 雪"},
    {code: 118, text: "晴れ のち 雨か雪"},
    {code: 119, text: "晴れ のち 雨か雷雨"},
    {code: 120, text: "晴れ 朝夕 一時 雨"},
    {code: 121, text: "晴れ 朝の内一時 雨"},
    {code: 122, text: "晴れ 夕方一時 雨"},
    {code: 123, text: "晴れ 山沿い 雷雨"},
    {code: 124, text: "晴れ 山沿い 雪"},
    {code: 125, text: "晴れ 午後は雷雨"},
    {code: 126, text: "晴れ 昼頃から雨"},
    {code: 127, text: "晴れ 夕方から雨"},
    {code: 128, text: "晴れ 夜は雨"},
    {code: 129, text: "晴れ 夜半から雨"},
    {code: 130, text: "朝の内 霧 後 晴れ"},
    {code: 131, text: "晴れ 明け方 霧"},
    {code: 132, text: "晴れ 朝夕 くもり"},
    {code: 140, text: "晴れ 時々 雨で雷を伴う"},
    {code: 160, text: "晴れ 一時 雪か雨"},
    {code: 170, text: "晴れ 時々 雪か雨"},
    {code: 181, text: "晴れ のち 雪か雨"},
    {code: 200, text: "くもり"},
    {code: 201, text: "くもり 時々 晴れ"},
    {code: 202, text: "くもり 一時 雨"},
    {code: 203, text: "くもり 時々 雨"},
    {code: 204, text: "くもり 一時 雪"},
    {code: 205, text: "くもり 時々 雪"},
    {code: 206, text: "くもり 一時 雨か雪"},
    {code: 207, text: "くもり 時々 雨か雪"},
    {code: 208, text: "くもり 一時 雨か雷雨"},
    {code: 209, text: "霧"},
    {code: 210, text: "くもり のち時々 晴れ"},
    {code: 211, text: "くもり のち 晴れ"},
    {code: 212, text: "くもり のち一時 雨"},
    {code: 213, text: "くもり のち時々 雨"},
    {code: 214, text: "くもり のち 雨"},
    {code: 215, text: "くもり のち一時 雪"},
    {code: 216, text: "くもり のち時々 雪"},
    {code: 217, text: "くもり のち 雪"},
    {code: 218, text: "くもり のち 雨か雪"},
    {code: 219, text: "くもり のち 雨か雷雨"},
    {code: 220, text: "くもり 朝夕一時 雨"},
    {code: 221, text: "くもり 朝の内一時 雨"},
    {code: 222, text: "くもり 夕方一時 雨"},
    {code: 223, text: "くもり 日中時々 晴れ"},
    {code: 224, text: "くもり 昼頃から雨"},
    {code: 225, text: "くもり 夕方から雨"},
    {code: 226, text: "くもり 夜は雨"},
    {code: 227, text: "くもり 夜半から雨"},
    {code: 228, text: "くもり 昼頃から雪"},
    {code: 229, text: "くもり 夕方から雪"},
    {code: 230, text: "くもり 夜は雪"},
    {code: 231, text: "くもり海上海岸は霧か霧雨"},
    {code: 240, text: "くもり 時々雨で 雷を伴う"},
    {code: 250, text: "くもり 時々雪で 雷を伴う"},
    {code: 260, text: "くもり 一時 雪か雨"},
    {code: 270, text: "くもり 時々 雪か雨"},
    {code: 281, text: "くもり のち 雪か雨"},
    {code: 300, text: "雨"},
    {code: 301, text: "雨 時々 晴れ"},
    {code: 302, text: "雨 時々 止む"},
    {code: 303, text: "雨 時々 雪"},
    {code: 304, text: "雨か雪"},
    {code: 306, text: "大雨"},
    {code: 307, text: "風雨共に強い"},
    {code: 308, text: "雨で暴風を伴う"},
    {code: 309, text: "雨 一時 雪"},
    {code: 311, text: "雨 のち 晴れ"},
    {code: 313, text: "雨 のち くもり"},
    {code: 314, text: "雨 のち時々 雪"},
    {code: 315, text: "雨 のち 雪"},
    {code: 316, text: "雨か雪 のち 晴れ"},
    {code: 317, text: "雨か雪 のち くもり"},
    {code: 320, text: "朝の内雨 のち 晴れ"},
    {code: 321, text: "朝の内雨 のち くもり"},
    {code: 322, text: "雨 朝晩一時 雪"},
    {code: 323, text: "雨 昼頃から 晴れ"},
    {code: 324, text: "雨 夕方から 晴れ"},
    {code: 325, text: "雨 夜は晴"},
    {code: 326, text: "雨 夕方から雪"},
    {code: 327, text: "雨 夜は雪"},
    {code: 328, text: "雨 一時強く降る"},
    {code: 329, text: "雨 一時 みぞれ"},
    {code: 340, text: "雪か雨"},
    {code: 350, text: "雨で雷を伴う"},
    {code: 361, text: "雪か雨 のち 晴れ"},
    {code: 371, text: "雪か雨 のち くもり"},
    {code: 400, text: "雪"},
    {code: 401, text: "雪 時々 晴れ"},
    {code: 402, text: "雪 時々止む"},
    {code: 403, text: "雪 時々 雨"},
    {code: 405, text: "大雪"},
    {code: 406, text: "風雪強い"},
    {code: 407, text: "暴風雪"},
    {code: 409, text: "雪 一時 雨"},
    {code: 411, text: "雪 のち 晴れ"},
    {code: 413, text: "雪 のち くもり"},
    {code: 414, text: "雪 のち 雨"},
    {code: 420, text: "朝の内雪 のち 晴れ"},
    {code: 421, text: "朝の内雪 のち くもり"},
    {code: 422, text: "雪 昼頃から雨"},
    {code: 423, text: "雪 夕方から雨"},
    {code: 424, text: "雪 夜半から雨"},
    {code: 425, text: "雪 一時強く降る"},
    {code: 426, text: "雪 のち みぞれ"},
    {code: 427, text: "雪 一時 みぞれ"},
    {code: 450, text: "雪で雷を伴う"}
  ];
  static abstractCodes: WeatherCode[] = [
    { code: 1, text: "晴れ"},
    { code: 2, text: "晴れと曇り"},
    { code: 3, text: "晴れと雨"},
    { code: 4, text: "雨"},
    { code: 5, text: "雪"}
  ];

  // codeからtextを取得
  getText(code: number) {
    const target = Weather.weatherCodes.find(e => e.code === code);
    return target ? target.text : null;
  }


}
