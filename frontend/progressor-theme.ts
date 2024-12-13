import type { CustomThemeConfig } from "@skeletonlabs/tw-plugin";

export const progressorTheme: CustomThemeConfig = {
  name: "progressor-theme",
  properties: {
    // =~= Theme Properties =~=
    "--theme-font-family-base": `"Montserrat", sans-serif`,
    "--theme-font-family-heading": `'Bebas Neue', sans-serif`,
    "--theme-font-color-base": "0 0 0",
    "--theme-font-color-dark": "255 255 255",
    "--theme-rounded-base": "2px",
    "--theme-rounded-container": "8px",
    "--theme-border-base": "2px",
    // =~= Theme On-X Colors =~=
    "--on-primary": "0 0 0",
    "--on-secondary": "0 0 0",
    "--on-tertiary": "0 0 0",
    "--on-success": "0 0 0",
    "--on-warning": "0 0 0",
    "--on-error": "255 255 255",
    "--on-surface": "0 0 0",
    // =~= Theme Colors  =~=
    // primary | #FAC638
    "--color-primary-50": "254 246 225", // #fef6e1
    "--color-primary-100": "254 244 215", // #fef4d7
    "--color-primary-200": "254 241 205", // #fef1cd
    "--color-primary-300": "253 232 175", // #fde8af
    "--color-primary-400": "252 215 116", // #fcd774
    "--color-primary-500": "250 198 56", // #FAC638
    "--color-primary-600": "225 178 50", // #e1b232
    "--color-primary-700": "188 149 42", // #bc952a
    "--color-primary-800": "150 119 34", // #967722
    "--color-primary-900": "123 97 27", // #7b611b
    // secondary | #76a93d
    "--color-secondary-50": "234 242 226", // #eaf2e2
    "--color-secondary-100": "228 238 216", // #e4eed8
    "--color-secondary-200": "221 234 207", // #ddeacf
    "--color-secondary-300": "200 221 177", // #c8ddb1
    "--color-secondary-400": "159 195 119", // #9fc377
    "--color-secondary-500": "118 169 61", // #76a93d
    "--color-secondary-600": "106 152 55", // #6a9837
    "--color-secondary-700": "89 127 46", // #597f2e
    "--color-secondary-800": "71 101 37", // #476525
    "--color-secondary-900": "58 83 30", // #3a531e
    // tertiary | #f7dc6f
    "--color-tertiary-50": "254 250 233", // #fefae9
    "--color-tertiary-100": "253 248 226", // #fdf8e2
    "--color-tertiary-200": "253 246 219", // #fdf6db
    "--color-tertiary-300": "252 241 197", // #fcf1c5
    "--color-tertiary-400": "249 231 154", // #f9e79a
    "--color-tertiary-500": "247 220 111", // #f7dc6f
    "--color-tertiary-600": "222 198 100", // #dec664
    "--color-tertiary-700": "185 165 83", // #b9a553
    "--color-tertiary-800": "148 132 67", // #948443
    "--color-tertiary-900": "121 108 54", // #796c36
    // success | #8bc34a
    "--color-success-50": "238 246 228", // #eef6e4
    "--color-success-100": "232 243 219", // #e8f3db
    "--color-success-200": "226 240 210", // #e2f0d2
    "--color-success-300": "209 231 183", // #d1e7b7
    "--color-success-400": "174 213 128", // #aed580
    "--color-success-500": "139 195 74", // #8bc34a
    "--color-success-600": "125 176 67", // #7db043
    "--color-success-700": "104 146 56", // #689238
    "--color-success-800": "83 117 44", // #53752c
    "--color-success-900": "68 96 36", // #446024
    // warning | #fab700
    "--color-warning-50": "254 244 217", // #fef4d9
    "--color-warning-100": "254 241 204", // #fef1cc
    "--color-warning-200": "254 237 191", // #feedbf
    "--color-warning-300": "253 226 153", // #fde299
    "--color-warning-400": "252 205 77", // #fccd4d
    "--color-warning-500": "250 183 0", // #fab700
    "--color-warning-600": "225 165 0", // #e1a500
    "--color-warning-700": "188 137 0", // #bc8900
    "--color-warning-800": "150 110 0", // #966e00
    "--color-warning-900": "123 90 0", // #7b5a00
    // error | #b30000
    "--color-error-50": "244 217 217", // #f4d9d9
    "--color-error-100": "240 204 204", // #f0cccc
    "--color-error-200": "236 191 191", // #ecbfbf
    "--color-error-300": "225 153 153", // #e19999
    "--color-error-400": "202 77 77", // #ca4d4d
    "--color-error-500": "179 0 0", // #b30000
    "--color-error-600": "161 0 0", // #a10000
    "--color-error-700": "134 0 0", // #860000
    "--color-error-800": "107 0 0", // #6b0000
    "--color-error-900": "88 0 0", // #580000
    // surface | #a1a1a1
    "--color-surface-50": "241 241 241", // #f1f1f1
    "--color-surface-100": "236 236 236", // #ececec
    "--color-surface-200": "232 232 232", // #e8e8e8
    "--color-surface-300": "217 217 217", // #d9d9d9
    "--color-surface-400": "189 189 189", // #bdbdbd
    "--color-surface-500": "161 161 161", // #a1a1a1
    "--color-surface-600": "145 145 145", // #919191
    "--color-surface-700": "121 121 121", // #797979
    "--color-surface-800": "97 97 97", // #616161
    "--color-surface-900": "79 79 79", // #4f4f4f
  },
};
