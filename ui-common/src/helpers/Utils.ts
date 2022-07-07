const IP_ADDRESS_REGEX = /([0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})(:([0-9]{1,5}))?/


export class Utils {
  static isIpAddress = (address: string): boolean => IP_ADDRESS_REGEX.test(address)
  static lineNumbersInString = (code: string): number => code.split("\n").length;

  static humanFileSize(bytes, si = false, dp = 1) {
    const thresh = si ? 1000 : 1024;

    if (Math.abs(bytes) < thresh) {
      return bytes + ' B';
    }

    const units = si
      ? ['kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
      : ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
    let u = -1;
    const r = 10 ** dp;

    do {
      bytes /= thresh;
      ++u;
    } while (Math.round(Math.abs(bytes) * r) / r >= thresh && u < units.length - 1);


    return bytes.toFixed(dp) + ' ' + units[u];
  }

  static padTo2Digits = (num) => {
    return String(num).padStart(2, '0');
  }

  static getHoursAndMinutes = (protocolTimeKey) => {
    const time = new Date(protocolTimeKey)
    const hoursAndMinutes = Utils.padTo2Digits(time.getHours()) + ':' + Utils.padTo2Digits(time.getMinutes());
    return hoursAndMinutes;
  }

  static formatDate = (date) => {
    let d = new Date(date),
        month = '' + (d.getMonth() + 1),
        day = '' + d.getDate(),
        year = d.getFullYear();
    const hoursAndMinutes = Utils.getHoursAndMinutes(date);
    if (month.length < 2) 
        month = '0' + month;
    if (day.length < 2) 
        day = '0' + day;
    const newDate = [year, month, day].join('-');
    return [hoursAndMinutes, newDate].join(' ');
}

  static creatUniqueObjArrayByProp = (objArray, prop) => {
    const map = new Map(objArray.map((item) => [item[prop], item])).values()
    return Array.from(map);
  }

  static isJson = (str) => {
    try {
      JSON.parse(str);
    } catch (e) {
      return false;
    }
    return true;
  }

}
