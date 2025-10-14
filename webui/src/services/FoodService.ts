import axios from "axios";
import type { IApiResponse } from "./Common";

export interface IFood {
  key: string;
  name: string;
  brand: string;
  cal100: number;
  prot100: number;
  fat100: number;
  carb100: number;
  comment: string;
}

export interface IFoodSet {
  food: IFood;
  isEdit: boolean;
}

export async function getFoodList() {
  const resp = await axios.get<IApiResponse>("/api/food");

  if (resp.data.error) {
    throw new Error(resp.data.error);
  }

  let respFood: IFood[] = [];
  for (let f of resp.data.data) {
    respFood.push({
      key: f.key,
      name: f.name,
      brand: f.brand,
      cal100: f.cal100,
      comment: f.comment,
    } as IFood);
  }

  return respFood;
}

export async function getFood(key: string) {
  const resp = await axios.get<IApiResponse>(`/api/food/${key}`);

  if (resp.data.error) {
    throw new Error(resp.data.error);
  }

  let f = resp.data.data;
  return {
    key: f.key,
    name: f.name,
    brand: f.brand,
    cal100: f.cal100,
    prot100: f.prot100,
    fat100: f.fat100,
    carb100: f.carb100,
    comment: f.comment,
  } as IFood;
}

export async function delFood(key: string) {
  const resp = await axios.delete<IApiResponse>(`/api/food/${key}`);
  if (resp.data.error) {
    throw new Error(resp.data.error);
  }
  return;
}

export async function setFood(foodSet: IFoodSet) {
  const resp = await axios.post<IApiResponse>("/api/food/set", foodSet);
  if (resp.data.error) {
    throw new Error(resp.data.error);
  }
  return;
}