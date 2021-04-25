
export default function setTitle(subtitle) {
  if (subtitle) {
    document.title = `jottr - ${subtitle}`
  } else {
    document.title = 'jottr'
  }
}
