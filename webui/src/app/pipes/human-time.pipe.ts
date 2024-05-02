import { Pipe, PipeTransform } from "@angular/core";

@Pipe({
  name: "humanTime",
  standalone: true,
})
export class HumanTimePipe implements PipeTransform {
  transform(value: Date | string): string {
    if (typeof value === "string") {
      value = new Date(value);
    }
    const seconds = Math.round((Date.now() - value.getTime()) / 1000);
    const suffix = seconds < 0 ? "from now" : "ago";
    const absSeconds = Math.abs(seconds);

    const times = [
      absSeconds / 60 / 60 / 24 / 365, // years
      absSeconds / 60 / 60 / 24 / 30, // months
      absSeconds / 60 / 60 / 24 / 7, // weeks
      absSeconds / 60 / 60 / 24, // days
      absSeconds / 60 / 60, // hours
      absSeconds / 60, // minutes
      absSeconds, // seconds
    ];
    const names = ["year", "month", "week", "day", "hour", "minute", "second"];

    for (let i = 0; i < names.length; i++) {
      const time = Math.floor(times[i]);
      let name = names[i];
      if (time > 1) name += "s";

      if (time >= 1) return time + " " + name + " " + suffix;
    }
    return "0 seconds " + suffix;
  }
}
