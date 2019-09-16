/** @format */

import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { HomeRouteComponent } from './routes/home/home.route';
import { LoginRouteComponent } from './routes/login/login.route';

const routes: Routes = [
  {
    path: '',
    component: HomeRouteComponent,
  },
  {
    path: 'login',
    component: LoginRouteComponent,
  },
  // {
  //   path: '**',
  //   redirectTo: '/',
  //   pathMatch: 'full',
  // },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
