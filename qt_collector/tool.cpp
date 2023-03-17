#include "tool.h"

QString ConvertString2CSV(QString& s)
{

   for (auto& c:s){
       if(c=='"' || c==','){
           //触发添加双引号
            QString r;
            r.push_back('"');
            for (auto& c:s){
                if(c=='"'){
                    r.push_back('"');
                }
                r.push_back(c);
            }
            r.push_back('"');
            return r;
       }
   }
   //遍历一遍无逗号和双引号，直接返回
   return s;
}
