import axios from 'axios'
import dotenv from 'dotenv'
import { Octokit } from 'octokit'

dotenv.config()

type BloodResp = {
  updateTime: string
  citys: {
    city: string
    'A型': string
    'B型': string
    'O型': string
    'AB型':string
  }[]
}

const bloodApiUrl = process.env.BLOOD_API_URL
const gist_id = process.env.GIST_ID
const token = process.env.TOKEN
const DISPLAY_TEXT: { [key: string]: string } = {
  less: "偏低",
  normal: "正常",
  lack: "急缺",
}
const displayTemplate = [
  '｜血型／城市｜',
  '｜Ａ        ｜',
  '｜Ｂ        ｜',
  '｜Ｏ        ｜',
  '｜ＡB       ｜',
]

const octokit = new Octokit({ auth: `token ${token}`})
;(async () => {
  const bloodResp = (await axios.get(bloodApiUrl)).data as BloodResp

  bloodResp.citys.forEach((el) => {
    Object.values(el).map((data, i) => {
      if (!i) displayTemplate[i] += data + '｜'
      else {
        displayTemplate[i] += DISPLAY_TEXT[data] + '｜'
      }
    }) 
  })

  const gist = await octokit.rest.gists.get({ gist_id })
  const filename = '🩸 ' + Object.keys(gist.data.files)[0] + ' ' + bloodResp.updateTime
  const content = displayTemplate.join('\n')
  await octokit.rest.gists.update({
    gist_id,
    files: {
      [filename]: {
        filename,
        content,
      },
    },
  })
})()