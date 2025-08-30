import { NgModule } from "@angular/core";
import { ConfigService } from "./config.service";
import { PluginsModule } from "../plugins/plugins.module";

@NgModule({
  imports: [PluginsModule],
  providers: [ConfigService],
})
export class ConfigModule {}
