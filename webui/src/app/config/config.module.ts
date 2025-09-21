import { NgModule } from "@angular/core";
import { PluginsModule } from "../plugins/plugins.module";
import { ConfigService } from "./config.service";

@NgModule({
  imports: [PluginsModule],
  providers: [ConfigService],
})
export class ConfigModule {}
