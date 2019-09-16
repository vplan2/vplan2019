/** @format */

import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HomeRouteComponent } from './routes/home/home.route';
import { LoginRouteComponent } from './routes/login/login.route';

@NgModule({
  declarations: [AppComponent, HomeRouteComponent, LoginRouteComponent],
  imports: [BrowserModule, AppRoutingModule],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
